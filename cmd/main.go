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
	if os.Args[1] == "install" {
		interval := config.LaunchAgentInterval
		if len(os.Args) == 3 {
			argsInterval, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Println("Invalid interval value - using default of 1200 seconds")
			}
			interval = argsInterval
		}
		launchagent.CreateLaunchAgent(deps, interval)
	} else if (os.Args[1] == "clean") {
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
		nc := notify.NewNotificationChannel("com.test.goPR")
		_, err = nc.Send("PR Changes", changeString)
		if err != nil {
			log.Println("Error sending notification: ", err)
		}
	}
}
