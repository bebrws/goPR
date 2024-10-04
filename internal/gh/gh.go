package gh

import (
	"fmt"

	"github.com/bebrws/goPR/internal/store"
)

func GetRepoState(client GitHubPullRequestsClient, cfg *store.Config) (*store.GHState, error) {
	state := store.GHState{}

	for _, repo := range cfg.Repos {
		allPrs, err := paginate(nil, GetPRPaginator(client, repo.Org, repo.Repo))
		if err != nil {
			return &state, err
		}
		prs := []store.PR{}
		for _, pr := range allPrs {
			revs := []store.PRReview{}
			fmt.Printf("PR %d. PR title: %s\nPR body: %s\n", *pr.Number, *pr.Title, *pr.Body)

			allRevs, err := paginate(nil, GetReviewPaginator(client, repo.Org, repo.Repo, *pr.Number))
			if err != nil {
				return &state, err
			}
			for j, rev := range allRevs {
				cmts := []store.PRReviewComment{}

				fmt.Printf("  %d. Review: %s from %s\n", j+1, rev.GetBody(), rev.User.GetLogin())

				allRevComments, err := paginate(nil, GetReviewCommentsPaginator(client, repo.Org, repo.Repo, *pr.Number, *rev.ID))
				if err != nil {
					return &state, err
				}

				for j, revComment := range allRevComments {
					cmts = append(cmts, store.PRReviewComment{
						ID:        revComment.GetID(),
						UpdatedAt: revComment.GetUpdatedAt().Time,
						Login:     revComment.User.GetLogin(),
						Body:      revComment.GetBody(),
					})
					fmt.Printf("      %d. Review Comment: %s from %s\n", j+1, revComment.GetBody(), revComment.User.GetLogin())
				}
				revs = append(revs, store.PRReview{
					ID:       rev.GetID(),
					Login:    rev.User.GetLogin(),
					Body:     rev.GetBody(),
					Comments: cmts,
				})
			}
			prs = append(prs, store.PR{
				Number:  *pr.Number,
				Body:    *pr.Body,
				Reviews: revs,
			})
		}
		state.RepoStates = append(state.RepoStates, store.RepoState{
			Name: repo.Repo,
			PRs:  prs,
		})
	}
	return &state, nil
}
