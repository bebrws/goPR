package main

import (
	"encoding/json"
	"log"
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
		log.Fatal("GH_USER is not set")
		ghUser = "" // Default value
	}

	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		log.Fatal("GH_TOKEN is not set")
		ghToken = "" // Default value
	}
	client := github.NewClient(nil).WithAuthToken(ghToken)
	if client == nil {
		log.Fatal("Failed to create client")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
		return
	}

	configFilePath := filepath.Join(homeDir, config.ConfigFileName)

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Error reading config file:", err)
		return
	}

	var cfg config.Config
	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
		return
	}

	newGHState, err := gh.GetRepoState(client, &cfg)

	if err != nil {
		log.Fatal("Failed to get repo state")
	}

	oldState := config.GHState{}

	stateFilePath := filepath.Join(homeDir, config.StateFileName)
	oldStateData, err := os.ReadFile(stateFilePath)
	if err != nil {
		// TODO: Refactor
		newStateData, err := json.Marshal(oldState)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err)
		}
		err = os.WriteFile(stateFilePath, newStateData, 0644)
		if err != nil {
			log.Fatal("Error writing state file:", err)
		}
	} else {
		err = json.Unmarshal(oldStateData, &oldState)
		if err != nil {
			log.Fatal("Error unmarshalling JSON:", err)
		}
	}

	diffs := config.CompareStates(oldState, *newGHState)
	changeString := strings.Join(diffs, "\n")
	log.Println("Changes: ", changeString)

	// Write the new state to the state file
	newStateData, err := json.Marshal(newGHState)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}
	err = os.WriteFile(stateFilePath, newStateData, 0644)
	if err != nil {
		log.Fatal("Error writing state file:", err)
	}
	if len(diffs) != 0 {
		nc := notify.NewNotificationChannel("com.bebrws.goPR")
		_, err = nc.Send("PR Changes", changeString)
		if err != nil {
			log.Println("Error sending notification: ", err)
		}
	}
}
