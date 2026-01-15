package githubmodule

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service         Service
	defaultUsername string
}

func NewController(s Service, defaultUsername string) *Controller {
	return &Controller{
		service:         s,
		defaultUsername: defaultUsername,
	}
}

// GetMyReposHandler godoc
// @Summary      List my GitHub repos (public + private)
// @Description  Mengambil semua repository milik akun pemilik token
// @Tags         github
// @Produce      json
// @Param        visibility  query   string  false  "Visibility: all|public|private"  default(all)
// @Success      200         {array} RepoDTO
// @Failure      500         {object} map[string]string
// @Router       /github/me/repos [get]
func (c *Controller) GetMyReposHandler(ctx *fiber.Ctx) error {
	visibility := ctx.Query("visibility", "all") // all, public, private

	repos, err := c.service.ListMyRepos(context.Background(), visibility)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch repos",
		})
	}

	return ctx.JSON(repos)
}

// ListPublicByUserHandler godoc
// @Summary      List public repos by username
// @Description  Mengambil repository publik berdasarkan username GitHub
// @Tags         github
// @Produce      json
// @Param        username  path      string  true  "GitHub username"
// @Success      200       {array}   RepoDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /github/{username}/repos [get]
func (c *Controller) ListPublicByUserHandler(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	if username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing username",
		})
	}

	repos, err := c.service.ListPublicRepoByUser(context.Background(), username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch repos",
		})
	}

	return ctx.JSON(repos)
}

// ContributionsHandler godoc
// @Summary      GitHub contributions calendar
// @Description  Mengambil data kontribusi GitHub (jumlah per tanggal) untuk 1 tahun terakhir
// @Tags         github
// @Produce      json
// @Param        username  path   string  true  "GitHub username"
// @Success      200       {object} githubmodule.ContributionResponse
// @Failure      400       {object} map[string]string
// @Failure      500       {object} map[string]string
// @Router       /github/{username}/contributions [get]
func (c *Controller) ContributionsHandler(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	if username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing username",
		})
	}

	yearStr := ctx.Query("year", "")
	var from, to time.Time

	if yearStr == "" {
		// Default: 1 tahun ke belakang hari ini
		to = time.Now()
		from = to.AddDate(-1, 0, 0)
	} else {
		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 2008 { // GitHub mulai 2008
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid year",
			})
		}
		// Range: 1 Jan sampai 31 Des Tahun yang dipilih (pakai timezone UTS)
		from = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		to = time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	data, err := c.service.FetchContributions(context.Background(), username, from, to)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch contributions",
		})
	}

	return ctx.JSON(data)
}
