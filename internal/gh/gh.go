package gh

import (
	"fmt"

	"github.com/bebrws/goPR/config"
	"github.com/google/go-github/v65/github"
)

func GetRepoState(client *github.Client, cfg *config.Config) (*config.GHState, error) {
	state := config.GHState{}

	for _, repo := range cfg.Repos {
		allPrs, err := paginate(GetPRPaginator(client, repo.Org, repo.Repo))
		if err != nil {
			fmt.Println("Error getting PRs", err)
		}
		prs := []config.PR{}
		for _, pr := range allPrs {
			revs := []config.PRReview{}
			fmt.Printf("PR %d. PR title: %s\nPR body: %s\n", *pr.Number, *pr.Title, *pr.Body)

			allRevs, err := paginate(GetReviewPaginator(client, repo.Org, repo.Repo, *pr.Number))
			if err != nil {
				fmt.Println("Error getting PR Reviews", err)
			}
			for j, rev := range allRevs {
				cmts := []config.PRReviewComment{}
				
				fmt.Printf("  %d. Review: %s from %s\n", j+1, rev.GetBody(), rev.User.GetLogin())

				allRevComments, err := paginate(GetReviewCommentsPaginator(client, repo.Org, repo.Repo, *pr.Number, *rev.ID))
				if err != nil {
					fmt.Println("Error getting PR Reviews", err)
				}

				for j, revComment := range allRevComments {
					cmts = append(cmts, config.PRReviewComment{
						Login: revComment.User.GetLogin(),
						Body:  revComment.GetBody(),
					})
					fmt.Printf("      %d. Review Comment: %s from %s\n", j+1, revComment.GetBody(), revComment.User.GetLogin())
				}
				revs = append(revs, config.PRReview{
					Login:    rev.User.GetLogin(),
					Body:     rev.GetBody(),
					Comments: cmts,
				})
			}
			prs = append(prs, config.PR{
				Number: *pr.Number,
				Body:   *pr.Body,
				Reviews: revs,
			})
		}
		state.RepoStates = append(state.RepoStates, config.RepoState{
			Name: repo.Repo,
			PRs:  prs,
		})
	}
	return &state, nil
}
