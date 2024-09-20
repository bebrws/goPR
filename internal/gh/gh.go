package gh

import (
	"fmt"

	"github.com/bebrws/goPR/config"
	"github.com/google/go-github/v65/github"
)

func GetRepoState(client *github.Client, cfg *config.Config) (*config.CurrentState, error) {
	state := config.CurrentState{}

	for _, repo := range cfg.Repos {
	allPrs, err := paginate(GetPRPaginator(client, repo.Org, repo.Repo))
	if err != nil {
		fmt.Println("Error getting PRs", err)
	}

	for _, pr := range allPrs {
		fmt.Printf("PR %d. PR title: %s\nPR body: %s\n", *pr.Number, *pr.Title, *pr.Body)

		allRevs, err := paginate(GetReviewPaginator(client, repo.Org, repo.Repo, *pr.Number))
		if err != nil {
			fmt.Println("Error getting PR Reviews", err)
		}
		for j, rev := range allRevs {
			fmt.Printf("  %d. Review: %s from %s\n", j+1, rev.GetBody(), rev.User.GetLogin())

			allRevComments, err := paginate(GetReviewCommentsPaginator(client, repo.Org, repo.Repo, *pr.Number, *rev.ID))
			if err != nil {
				fmt.Println("Error getting PR Reviews", err)
			}

			for j, revComment := range allRevComments {
				fmt.Printf("      %d. Review Comment: %s from %s\n", j+1, revComment.GetBody(), revComment.User.GetLogin())
			}
		}
	}
}
return &state, nil
}