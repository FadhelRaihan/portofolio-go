package githubmodule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

type Service interface {
	ListPublicRepoByUser(ctx context.Context, username string) ([]RepoDTO, error)
	ListMyRepos(ctx context.Context, visibility string) ([]RepoDTO, error)
	FetchContributions(ctx context.Context, username string, from, to time.Time) (*ContributionResponse, error)
}

type service struct {
	client     *github.Client
	httpClient *http.Client
	token      string
}

type RepoDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Private     bool   `json:"private"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

const githubGraphQLEndpoint = "https://api.github.com/graphql"

type ContributionDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ContributionResponse struct {
	Days []ContributionDay `json:"days"`
}

func NewService(token string) Service {
	ctx := context.Background()

	var httpClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(ctx, ts)
	} else {
		httpClient = http.DefaultClient
	}

	client := github.NewClient(httpClient)

	return &service{
		client:     client,
		httpClient: httpClient,
		token:      token,
	}
}

// Public repos user
func (s *service) ListPublicRepoByUser(ctx context.Context, username string) ([]RepoDTO, error) {
	opts := &github.RepositoryListByUserOptions{
		Type: "public",
		Sort: "updated",
	}

	repos, _, err := s.client.Repositories.ListByUser(ctx, username, opts)
	if err != nil {
		return nil, err
	}

	return mapRepos(repos), nil
}

// Semua repos milik akun yang ter auth (public + private)
func (s *service) ListMyRepos(ctx context.Context, visibility string) ([]RepoDTO, error) {
	// Visibility: "all" | "public" | "private"
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Visibility: "visibility",
		Sort:       "updated",
	}

	repos, _, err := s.client.Repositories.ListByAuthenticatedUser(ctx, opts)
	if err != nil {
		return nil, err
	}

	return mapRepos(repos), nil
}

func mapRepos(repos []*github.Repository) []RepoDTO {
	res := make([]RepoDTO, 0, len(repos))
	for _, r := range repos {
		res = append(res, RepoDTO{
			Name:        r.GetName(),
			Description: r.GetDescription(),
			HTMLURL:     r.GetHTMLURL(),
			Private:     r.GetPrivate(),
			CreatedAt:   r.GetCreatedAt().Format("2006-01-02"),
			UpdatedAt:   r.GetUpdatedAt().Format("2006-01-02"),
		})
	}
	return res
}

// Contributions
func (s *service) FetchContributions(ctx context.Context, username string, from, to time.Time) (*ContributionResponse, error) {
	query := `
	query($login: String!, $from: DateTime!, $to: DateTime!) {
		user(login: $login) {
			contributionsCollection(from: $from, to: $to) {
				contributionCalendar {
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"login": username,
		"from":  from.Format(time.RFC3339),
		"to":    to.Format(time.RFC3339),
	}

	bodyBytes, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubGraphQLEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github graphql status: %s body=%s", resp.Status, string(b))
	}

	var raw struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						Weeks []struct {
							ContributionDay []struct {
								Date              string `json:"date"`
								ContributionCount int    `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	days := make([]ContributionDay, 0, 370)
	for _, w := range raw.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, d := range w.ContributionDay {
			days = append(days, ContributionDay{
				Date:  d.Date,
				Count: d.ContributionCount,
			})
		}
	}

	return &ContributionResponse{Days: days}, nil
}
