package di

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/store"
	"github.com/google/go-github/v65/github"
)

type Deps struct {
	HomeDir string
	ExecutablePath string
	Client         gh.GitHubPullRequestsClient
	Config         store.Config
	OldState       store.GHState
	StateFilePath  string
	ConfigFilePath string
}

func GetGHTokenOrPanic() string {
	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		log.Fatal("GH_TOKEN is not set")
		ghToken = "" // Default value
	}
	return ghToken
}

func GetGHPPRClientOrPanic(ghToken string) *github.PullRequestsService {
	if ghToken == "" {
		log.Fatal("GH_TOKEN is not set")
	}
	return gh.NewPRClientOrPanic(ghToken)
}

func GetHomeDirOrPanic() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}
	return homeDir
}

func GetCfgOrPanic(configFilePath string) store.Config {
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}
	var cfg store.Config
	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	return cfg
}

func GetOldStateOrPanic(stateFilePath string) store.GHState {
	oldState := store.GHState{}
	oldStateData, err := os.ReadFile(stateFilePath)
	if err != nil {
		log.Println("Error reading state oldState file (will attempt to create one):", err)
		store.WriteState(stateFilePath, &oldState)
	} else {
		err = json.Unmarshal(oldStateData, &oldState)
		if err != nil {
			log.Fatal("Error unmarshalling oldState JSON:", err)
		}
	}
	return oldState
}

func NewDepsOrPanic() *Deps {
	homeDir := GetHomeDirOrPanic()
	ghToken := GetGHTokenOrPanic()
	client := GetGHPPRClientOrPanic(ghToken)
	configFilePath := filepath.Join(homeDir, config.ConfigFileName)
	stateFilePath := filepath.Join(homeDir, config.StateFileName)
	cfg := GetCfgOrPanic(configFilePath)
	oldState := GetOldStateOrPanic(stateFilePath)
	return &Deps{
		HomeDir: homeDir,
		ExecutablePath: os.Args[0],
		Client:         client,
		Config:         cfg,
		OldState:       oldState,
		StateFilePath:  stateFilePath,
		ConfigFilePath: configFilePath,
	}
}
