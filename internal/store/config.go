package store

import (
	"os"

	"github.com/bebrws/goPR/internal/models"
	"github.com/google/go-github/v65/github"
)

func CreateExConfig(repos []*github.Repository) *models.Config {
	exConfig := &models.Config{
		GHToken: os.Getenv("GITHUB_TOKEN"),
		Repos:   []models.Repo{},
	}
	for _, repo := range repos {
		exConfig.Repos = append(exConfig.Repos, models.Repo{
			Org:  repo.GetOwner().GetLogin(),
			Repo: repo.GetName(),
		})
	}
	return exConfig
}
