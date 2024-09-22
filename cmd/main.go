package main

import (
	"log"
	"strings"

	"github.com/bebrws/goPR/internal/di"
	"github.com/bebrws/goPR/internal/gh"
	"github.com/bebrws/goPR/internal/launchagent"
	"github.com/bebrws/goPR/internal/store"

	"github.com/bebrws/goPR/internal/notify"
)

func main() {
	nc := notify.NewNotificationChannel("com.bebrws.goPR")
	_, err := nc.Send("PR Changes", "changeString")
	if err != nil {
		log.Println("Error sending notification: ", err)
	}

	deps := di.NewDepsOrPanic()

	launchagent.CreateLaunchAgent(deps)
	
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
		nc := notify.NewNotificationChannel("com.bebrws.goPR")
		_, err = nc.Send("PR Changes", changeString)
		if err != nil {
			log.Println("Error sending notification: ", err)
		}
	}
}
