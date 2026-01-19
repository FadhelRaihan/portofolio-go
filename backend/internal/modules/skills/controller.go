package skills

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	MaxFileSize   = 5 * 1024 * 1024 // 5MB
	StorageBucket = "skills"
	UploadTimeout = 30 * time.Second
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
}

type Controller struct {
	service      Service
	supabaseURL  string
	supabaseKey  string
	httpClient   *http.Client
}

func NewController(service Service, supabaseURL, supabaseKey string) *Controller {
	return &Controller{
		service:     service,
		supabaseURL: supabaseURL,
		supabaseKey: supabaseKey,
		httpClient: &http.Client{
			Timeout: UploadTimeout,
		},
	}
}

// Validate file upload
func (c *Controller) validateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("invalid file type. Allowed: jpg, jpeg, png, gif, webp, svg")
	}

	return nil
}

// Upload file to Supabase Storage using REST API
func (c *Controller) uploadFile(file *multipart.FileHeader) (string, error) {
	if err := c.validateFile(file); err != nil {
		return "", err
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueName := uuid.New().String() + ext

	// Prepare upload URL
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		c.supabaseURL, StorageBucket, uniqueName)

	// Create request
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.supabaseKey)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))
	req.Header.Set("apikey", c.supabaseKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Return public URL
	imageURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		c.supabaseURL, StorageBucket, uniqueName)

	return imageURL, nil
}

// Delete file from Supabase Storage using REST API
func (c *Controller) deleteFile(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Extract filename from URL
	parts := strings.Split(imageURL, "/")
	if len(parts) == 0 {
		return nil
	}
	filename := parts[len(parts)-1]

	// Prepare delete URL
	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		c.supabaseURL, StorageBucket, filename)

	// Create request
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		log.Printf("Warning: failed to create delete request: %v", err)
		return nil // Don't fail the operation
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.supabaseKey)
	req.Header.Set("apikey", c.supabaseKey)

	// Execute request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Warning: failed to delete file %s: %v", filename, err)
		return nil // Don't fail the operation
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Warning: delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	var input models.CreateSkillInput

	// Parse multipart form
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Handle file upload (optional)
	var imageURL string
	file, err := ctx.FormFile("image")
	if err == nil && file != nil {
		imageURL, err = c.uploadFile(file)
		if err != nil {
			log.Printf("Upload error: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Create skill
	skill, err := c.service.Create(context.Background(), &input, imageURL)
	if err != nil {
		log.Printf("Create skill error: %v", err)
		// If creation fails, delete uploaded image
		if imageURL != "" {
			c.deleteFile(imageURL)
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create skill",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(skill)
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Skill ID is required",
		})
	}

	var input models.UpdateSkillInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get existing skill to get old image URL
	existing, err := c.service.GetByID(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Skill not found",
		})
	}

	oldImageURL := existing.Image

	// Handle new file upload (optional)
	var imageURL string
	file, err := ctx.FormFile("image")
	if err == nil && file != nil {
		imageURL, err = c.uploadFile(file)
		if err != nil {
			log.Printf("Upload error: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Update skill
	skill, err := c.service.Update(context.Background(), id, &input, imageURL)
	if err != nil {
		log.Printf("Update skill error: %v", err)
		// If update fails, delete newly uploaded image
		if imageURL != "" {
			c.deleteFile(imageURL)
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update skill",
		})
	}

	// Delete old image if new one was uploaded
	if imageURL != "" && oldImageURL != "" {
		c.deleteFile(oldImageURL)
	}

	return ctx.JSON(skill)
}

func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Skill ID is required",
		})
	}

	skill, err := c.service.GetByID(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Skill not found",
		})
	}

	return ctx.JSON(skill)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Skill ID is required",
		})
	}

	// Delete skill and get old image URL
	oldImageURL, err := c.service.Delete(context.Background(), id)
	if err != nil {
		log.Printf("Delete skill error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete skill",
		})
	}

	// Delete image from storage
	if oldImageURL != "" {
		c.deleteFile(oldImageURL)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Skill deleted successfully",
	})
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	skills, err := c.service.List(context.Background())
	if err != nil {
		log.Printf("List skills error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list skills",
		})
	}

	return ctx.JSON(skills)
}