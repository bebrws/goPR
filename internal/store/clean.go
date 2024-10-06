package store

import (
	"fmt"
	"os"
)

func CleanupStateAndConfig(stateFilePath string, configFilePath string) {
	// Remove the state and config files
	err := os.Remove(stateFilePath)
	if err != nil {
		fmt.Printf("Error removing state file: %v\n", err)
	}
	err = os.Remove(configFilePath)
	if err != nil {
		fmt.Printf("Error removing config file: %v\n", err)
	}
} 