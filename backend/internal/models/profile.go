package models

import "time"

type Profile struct {
	ID          string    `db:"id"`
	FullName    string    `db:"full_name"`
	Headline    string    `db:"headline"`
	Bio         string    `db:"bio"`
	Location    string    `db:"location"`
	GithubURL   string    `db:"github_url"`
	LinkedinURL string    `db:"linkedin_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// DTO untuk response
type ProfileDTO struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	Headline    string `json:"headline"`
	Bio         string `json:"bio"`
	Location    string `json:"location,omitempty"`
	GithubURL   string `json:"github_url,omitempty"`
	LinkedinURL string `json:"linkedin_url,omitempty"`
}
