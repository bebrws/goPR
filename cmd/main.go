package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"path/filepath"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/notify"
	"github.com/google/go-github/v65/github"
)

func main() {	
	ghUser := os.Getenv("GH_USER")
	if ghUser == "" {
		fmt.Println("GH_USER is not set")
		ghUser = "" // Default value
	}

	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		fmt.Println("GH_TOKEN is not set")
		ghToken = "" // Default value
	}
	client := github.NewClient(nil).WithAuthToken(ghToken)
	if client == nil {
		fmt.Println("Failed to create client")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configFilePath := filepath.Join(homeDir, config.ConfigFileName)
	
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	
	var cfg config.Config
	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	
	newGHState, err := gh.GetRepoState(client, &cfg)

	if err != nil {
		fmt.Println("Failed to get repo state")
		os.Exit(1)
	}

	stateFilePath := filepath.Join(homeDir, config.StateFileName)
	oldStateData, err := os.ReadFile(stateFilePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var oldState config.GHState
	err = json.Unmarshal(oldStateData, &oldState)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	diffs := config.CompareStates(oldState, *newGHState)
	changeString := strings.Join(diffs, "\n")
	fmt.Println("Changes: ", changeString)

	nc := notify.NewNotificationChannel("com.bebrws.goPR")
	_, err = nc.Send("PR Changes", changeString)
	if err != nil {
		fmt.Println(err)
	}
}
