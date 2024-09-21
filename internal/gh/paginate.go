package gh

import (
	"context"
	"log"

	"github.com/bebrws/goPR/config"
	"github.com/google/go-github/v65/github"
)

type RateLimitedPage struct {
	github.ListOptions
	// Rate *github.Rate
	resp *github.Response
}

func NewRateLimitedPage(opts *RateLimitedPage, resp *github.Response) *RateLimitedPage {
	return &RateLimitedPage{
		ListOptions: opts.ListOptions,
		resp: resp,
	}
}

func (lo *RateLimitedPage) Update(resp *github.Response) {
	lo.resp = resp
	lo.Page = lo.resp.NextPage
}

func (lo *RateLimitedPage) GetRateLimitRemaining() int {
	return lo.resp.Rate.Remaining
}

// What is a RepositoryComment??
// https://github.com/google/go-github/blob/23592519310b534946cbb8527ad6a92d5d29ddd0/github/event_types.go#L71C6-L71C24
// github.RepositoryComment Is inside CommitCommentEvent 
// which `CommitCommentEvent is triggered when a commit comment is created.`

type PaginateAbleGithubTypes interface {
	*github.PullRequest | *github.Comment | *github.PullRequestComment | *github.PullRequestReview | *github.RepositoryComment
}

type GitHubPullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	ListReviews(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
	ListReviewComments(ctx context.Context, owner, repo string, number int, reviewID int64, opts *github.ListOptions) ([]*github.PullRequestComment, *github.Response, error)
}

func paginate[R PaginateAbleGithubTypes](opts *RateLimitedPage, pf func(opts *RateLimitedPage) ([]R, *github.Response, error)) ([]R, error) {
	var listOps *RateLimitedPage
	if opts == nil {
		rlp := RateLimitedPage{
			ListOptions: github.ListOptions{
				PerPage: config.PerPage,
			},
		}
		listOps = NewRateLimitedPage(&rlp, nil)
	} else {
		listOps = NewRateLimitedPage(opts, nil)
	}
	allItems := []R{}
	for {
		items, resp, err := pf(listOps)
		listOps.Update( resp )

		if resp.Rate.Remaining == 0 {
			log.Fatal("Rate limit reached, time to panic") // TODO: Return a special error or boolean?
		} else if resp != nil {
			log.Println("Rate limit remaining:", resp.Rate.Remaining)
		}
		
		if err != nil {	
			return allItems, err
		}
		if len(items) == 0 {
			break
		}
		allItems = append(allItems, items...)
		if len(items) < config.PerPage {
			break
		}
	}
	return allItems, nil
}

func GetPRPaginator(client GitHubPullRequestsClient, org, repo string) func(opts *RateLimitedPage) ([]*github.PullRequest, *github.Response, error) {
	return func(opts *RateLimitedPage) ([]*github.PullRequest, *github.Response, error) {
		prLO := github.PullRequestListOptions{
			ListOptions: opts.ListOptions,
			State: config.PrState,
		}
		return client.List(context.Background(), org, repo, &prLO)
	}
}

func GetReviewPaginator(client GitHubPullRequestsClient, org, repo string, prNumber int) func(opts *RateLimitedPage) ([]*github.PullRequestReview, *github.Response, error) {
	return func(opts *RateLimitedPage) ([]*github.PullRequestReview, *github.Response, error) {
		return client.ListReviews(context.Background(), org, repo, prNumber, &opts.ListOptions)
	}
}

func GetReviewCommentsPaginator(client GitHubPullRequestsClient, org, repo string, prNumber int, revID int64) func(opts *RateLimitedPage) ([]*github.PullRequestComment, *github.Response, error) {
	return func(opts *RateLimitedPage) ([]*github.PullRequestComment, *github.Response, error) {
		return client.ListReviewComments(context.Background(), org, repo, prNumber, revID, &opts.ListOptions)
	}
}