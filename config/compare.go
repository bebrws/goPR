package config

import (
	"log"
	"sort"
	"strings"
)

// Compare function to detect differences between two GHState structs
func CompareStates(oldState, newState GHState) []string {
	var changes []string

	// TODO: Need to "deep" sort?
	sort.Slice(oldState.RepoStates, func(i, j int) bool {
		return oldState.RepoStates[i].PRs[0].Number < oldState.RepoStates[j].PRs[0].Number
	})
	sort.Slice(newState.RepoStates, func(i, j int) bool {
		return newState.RepoStates[i].PRs[0].Number < newState.RepoStates[j].PRs[0].Number
	})
	

	log.Println("my diff: ", strings.Join(changes, ""))

	return changes
}
