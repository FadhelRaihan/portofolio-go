package models

import "time"

type Skill struct {
    ID         string    `db:"id" json:"id"`
    Name       string    `db:"name" json:"name"`
    Experience int       `db:"experience" json:"experience"`
    Image      string    `db:"image" json:"image"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type SkillDTO struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    Experience int    `json:"experience"`
    Image      string `json:"image,omitempty"`
}

// Input validation
type CreateSkillInput struct {
    Name       string `json:"name" form:"name" validate:"required,min=2,max=100"`
    Experience int    `json:"experience" form:"experience" validate:"required,min=0,max=100"`
}

type UpdateSkillInput struct {
    Name       *string `json:"name,omitempty" form:"name" validate:"omitempty,min=2,max=100"`
    Experience *int    `json:"experience,omitempty" form:"experience" validate:"omitempty,min=0,max=100"`
}