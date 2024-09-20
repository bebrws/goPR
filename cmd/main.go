package main

import (
	"fmt"
	"os"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/notify"
	"github.com/google/go-github/v65/github"
)

func main() {
	nc := notify.NewNotificationChannel("com.bebrws.goPR")
	_, err := nc.Send("Hello", "World")
	if err != nil {
		fmt.Println(err)
	}
	
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

	cfg := config.Config{
	}

	state, err := gh.GetRepoState(client, &cfg)
	if err != nil {
		fmt.Println("Failed to get repo state")
		os.Exit(1)
	}

	fmt.Println(state)
}
