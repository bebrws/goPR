package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/di"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/launchagent"
	"github.com/bebrws/goPR/internal/notify"
	"github.com/bebrws/goPR/internal/store"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	deps, err := di.NewDeps()
	if err != nil {
		logrus.Fatalf("Failed to initialize dependencies: %v", err)
	}

	printUsage(deps.ExecutablePath)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			handleInstall(deps)
		case "clean":
			handleClean(deps)
		default:
			logrus.Fatalf("Unknown command: %s", os.Args[1])
		}
	}

	newGHState, err := gh.GetRepoState(deps.Client, &deps.Config)
	if err != nil {
		handleInitGetStateError(err)
	}

	store.WriteState(deps.StateFilePath, newGHState)

	diffs := store.CompareStates(deps.OldState, *newGHState)
	changeString := strings.Join(diffs, "\n")
	logrus.Info("Changes: ", changeString)

	if len(diffs) != 0 {
		sendNotification(changeString)
	}
}

func printUsage(executablePath string) {
	fmt.Printf("goPR - Located at %s takes the following arguments:\n", executablePath)
	fmt.Printf("%s install - to install the launch agent\n", executablePath)
	fmt.Printf("%s clean - to clean the launch agent and state and config file\n", executablePath)
}

func handleInstall(deps *di.Deps) {
	interval := config.LaunchAgentInterval
	if len(os.Args) == 3 {
		argsInterval, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logrus.Warn("Invalid interval value - using default of 1200 seconds")
		} else {
			interval = argsInterval
		}
	}
	logrus.Info("Installing launch agent - interval: ", interval)
	launchagent.CreateLaunchAgent(deps, interval)
}

func handleClean(deps *di.Deps) {
	logrus.Info("Deleting launch agent and state and config file")
	launchagent.CleanLaunchAgent(deps)
	store.CleanupStateAndConfig(deps.StateFilePath, deps.ConfigFilePath)
}

func handleInitGetStateError(err error) {
	if rateLimitErr, ok := err.(*gh.RateLimitError); ok {
		logrus.Fatal("Rate limit reached: ", rateLimitErr)
	}
	logrus.Warn("Failed to get repo state: ", err)
}

func sendNotification(changeString string) {
	logrus.Info("Sending notification")
	nc := notify.NewNotificationChannel(config.BundleID)
	_, err := nc.Send("PR Changes", changeString)
	if err != nil {
		logrus.Error("Error sending notification: ", err)
	}
}
