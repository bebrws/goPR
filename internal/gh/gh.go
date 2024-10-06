package gh

import (
	"github.com/bebrws/goPR/internal/models"
	"github.com/sirupsen/logrus"
)

func GetRepoState(client GitHubPullRequestsClient, cfg *models.Config) (*models.GHState, error) {
	state := models.GHState{}

	for _, repo := range cfg.Repos {
		allPrs, err := Paginate(nil, GetPRPaginator(client, repo.Org, repo.Repo))
		if err != nil {
			return &state, err
		}
		prs := []models.PR{}
		for _, pr := range allPrs {
			revs := []models.PRReview{}
			logrus.Infof("PR %d. PR title: %s\nPR body: %s\n", *pr.Number, *pr.Title, *pr.Body)

			allRevs, err := Paginate(nil, GetReviewPaginator(client, repo.Org, repo.Repo, *pr.Number))
			if err != nil {
				return &state, err
			}
			for j, rev := range allRevs {
				cmts := []models.PRReviewComment{}

				logrus.Infof("  %d. Review: %s from %s\n", j+1, rev.GetBody(), rev.User.GetLogin())

				allRevComments, err := Paginate(nil, GetReviewCommentsPaginator(client, repo.Org, repo.Repo, *pr.Number, *rev.ID))
				if err != nil {
					return &state, err
				}

				for j, revComment := range allRevComments {
					cmts = append(cmts, models.PRReviewComment{
						ID:        revComment.GetID(),
						UpdatedAt: revComment.GetUpdatedAt().Time,
						Login:     revComment.User.GetLogin(),
						Body:      revComment.GetBody(),
					})
					logrus.Infof("      %d. Review Comment: %s from %s\n", j+1, revComment.GetBody(), revComment.User.GetLogin())
				}
				revs = append(revs, models.PRReview{
					ID:       rev.GetID(),
					Login:    rev.User.GetLogin(),
					Body:     rev.GetBody(),
					Comments: cmts,
				})
			}
			prs = append(prs, models.PR{
				Number:  *pr.Number,
				Body:    *pr.Body,
				Reviews: revs,
			})
		}
		state.RepoStates = append(state.RepoStates, models.RepoState{
			Name: repo.Repo,
			PRs:  prs,
		})
	}
	return &state, nil
}
