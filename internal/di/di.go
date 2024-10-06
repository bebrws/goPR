package di

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/models"
	"github.com/bebrws/goPR/internal/store"
	"github.com/google/go-github/v65/github"
	"github.com/sirupsen/logrus"
)

type Deps struct {
	HomeDir        string
	ExecutablePath string
	Client         gh.GitHubPullRequestsClient
	Config         models.Config
	OldState       models.GHState
	StateFilePath  string
	ConfigFilePath string
}

func GetGHPPRClient(ghToken string) (*github.PullRequestsService, error) {
	if ghToken == "" {
		return nil, errors.New("gh token is empty")
	}
	return gh.NewPRClient(ghToken), nil
}

func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func GetOrCreateDefaultConfig(configFilePath string) (models.Config, error) {
	configData, err := os.ReadFile(configFilePath)
	var cfg models.Config
	if err != nil {
		logrus.Warn("Config file not found, creating default config")
		client := gh.NewGHClient(os.Getenv("GITHUB_TOKEN"))
		gh_username, _, err := client.Users.Get(context.Background(), "")
		if err != nil {
			cfg.GHToken = os.Getenv("GITHUB_TOKEN")
			return cfg, nil
		}
		repos, err := gh.Paginate(nil, gh.GetRepoPaginator(client.Repositories, gh_username.GetLogin()))
		if err != nil {
			cfg.GHToken = os.Getenv("GITHUB_TOKEN")
			return cfg, err
		}
		cfg = *store.CreateExConfig(repos)
		logrus.Info("Created default config with all repos. Edit config to select specific repos at ", configFilePath)
		configFile, err := os.Create(configFilePath)
		if err != nil {
			logrus.Warn("Failed to create config file: ", err)
			return cfg, err
		}
		defer configFile.Close()

		encoder := json.NewEncoder(configFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(cfg); err != nil {
			logrus.Warn("Failed to write config to file: ", err)
			return cfg, err
		}
		return cfg, nil
	}

	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		return models.Config{}, err
	}
	return cfg, nil
}

func GetOldState(stateFilePath string) (models.GHState, error) {
	oldState := models.GHState{}
	oldStateData, err := os.ReadFile(stateFilePath)
	if err != nil {
		store.WriteState(stateFilePath, &oldState)
		return oldState, err
	}
	err = json.Unmarshal(oldStateData, &oldState)
	if err != nil {
		return oldState, err
	}
	return oldState, nil
}

func NewDeps() (*Deps, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return nil, err
	}
	configFilePath := filepath.Join(homeDir, config.ConfigFileName)
	stateFilePath := filepath.Join(homeDir, config.StateFileName)
	cfg, err := GetOrCreateDefaultConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	oldState, err := GetOldState(stateFilePath)
	if err != nil {
		return nil, err
	}
	client, err := GetGHPPRClient(cfg.GHToken)
	if err != nil {
		return nil, err
	}
	return &Deps{
		HomeDir:        homeDir,
		ExecutablePath: os.Args[0],
		Client:         client,
		Config:         cfg,
		OldState:       oldState,
		StateFilePath:  stateFilePath,
		ConfigFilePath: configFilePath,
	}, nil
}
