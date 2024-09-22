package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/di"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/launchagent"
	"github.com/bebrws/goPR/internal/store"

	"github.com/bebrws/goPR/internal/notify"
)

func main() {
	deps := di.NewDepsOrPanic()

	fmt.Printf("goPR - Located at %s takes the following arguments:\n", deps.ExecutablePath)
	fmt.Printf("%s install - to install the launch agent\n", deps.ExecutablePath)
	fmt.Printf("%s clean - to clean the launch agent and state and config file\n", deps.ExecutablePath)
	if len(os.Args) > 1 && os.Args[1] == "install" {
		interval := config.LaunchAgentInterval
		if len(os.Args) == 3 {
			argsInterval, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Println("Invalid interval value - using default of 1200 seconds")
			}
			interval = argsInterval
		}
		fmt.Println("Instaling launch agent - interval: ", interval)
		launchagent.CreateLaunchAgent(deps, interval)
	} else if len(os.Args) > 1 && os.Args[1] == "clean" {
		fmt.Println("Deleting launch agent and state and config file")
		launchagent.CleanLaunchAgent(deps)
		store.CleanupStateAndConfig(deps.StateFilePath, deps.ConfigFilePath)
	}
	
	newGHState, err := gh.GetRepoState(deps.Client, &deps.Config)
	if err != nil {
		if err.(*gh.RateLimitError) != nil {
			log.Fatal("Rate limit reached")
		}
		log.Fatal("Failed to get repo state: ", err)
	}

	store.WriteState(deps.StateFilePath, newGHState)

	diffs := store.CompareStates(deps.OldState, *newGHState)
	changeString := strings.Join(diffs, "\n")
	log.Println("Changes: ", changeString)
	
	if len(diffs) != 0 {
		log.Println("Sending notification")
		nc := notify.NewNotificationChannel(config.BundleID)
		_, err = nc.Send("PR Changes", changeString)
		if err != nil {
			log.Println("Error sending notification: ", err)
		}
	}
}
