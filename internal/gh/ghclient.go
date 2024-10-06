package gh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bebrws/goPR/config"
	"github.com/google/go-github/v65/github"
	"github.com/sirupsen/logrus"
)

type RateLimitError struct {
	Code    int
	Message string
}

// Implement the Error() method to satisfy the error interface
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewRateLimitError(message string) *RateLimitError {
	return &RateLimitError{
		Code:    http.StatusTooManyRequests,
		Message: message,
	}
}

type RateLimitedPage struct {
	github.ListOptions
	resp *github.Response
}

func NewRateLimitedPage(opts *RateLimitedPage, resp *github.Response) *RateLimitedPage {
	return &RateLimitedPage{
		ListOptions: opts.ListOptions,
		resp:        resp,
	}
}

func (lo *RateLimitedPage) Update(resp *github.Response) {
	lo.resp = resp
	lo.Page = lo.resp.NextPage
}

func (lo *RateLimitedPage) GetRateLimitRemaining() int {
	return lo.resp.Rate.Remaining
}

type PaginateAbleGithubTypes interface {
	*github.PullRequest | *github.Comment | *github.PullRequestComment | *github.PullRequestReview | *github.RepositoryComment | *github.Repository
}

type GitHubPullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	ListReviews(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
	ListReviewComments(ctx context.Context, owner, repo string, number int, reviewID int64, opts *github.ListOptions) ([]*github.PullRequestComment, *github.Response, error)
}

type GitHubRepositoryClient interface {
	List(ctx context.Context, owner string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
}

func NewGHClient(ghToken string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(ghToken)
	if client == nil {
		logrus.Fatal("Failed to create GitHub client")
	}
	return client
}

func NewRepoClient(ghToken string) *github.RepositoriesService {
	client := NewGHClient(ghToken)
	if client.PullRequests == nil {
		logrus.Fatal("Failed to get GitHub PullRequest client")
	}
	return client.Repositories
}

func NewPRClient(ghToken string) *github.PullRequestsService {
	client := NewGHClient(ghToken)
	if client.PullRequests == nil {
		logrus.Fatal("Failed to get GitHub PullRequest client")
	}
	return client.PullRequests
}

func Paginate[R PaginateAbleGithubTypes](opts *RateLimitedPage, pf func(opts *RateLimitedPage) ([]R, *github.Response, error)) ([]R, error) {
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
		listOps.Update(resp)

		if resp.Rate.Remaining <= 0 {
			logrus.Warn("Rate limit reached, time to panic")
			return allItems, NewRateLimitError("Rate limit reached")
		} else if resp != nil {
			logrus.Info("Rate limit remaining:", resp.Rate.Remaining)
		}

		if err != nil {
			// TODO: Its OK if we cannot access some repos, but then we will miss out on updates from the PRs in repos
			// that would have been returned in this pagination..
			if _, ok := any(items).([]*github.Repository); ok {
				logrus.Info("Continuing pagination for repositories")
				continue
			}
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

func GetRepoPaginator(client GitHubRepositoryClient, org string) func(opts *RateLimitedPage) ([]*github.Repository, *github.Response, error) {
	return func(opts *RateLimitedPage) ([]*github.Repository, *github.Response, error) {
		prLO := github.RepositoryListOptions{
			ListOptions: opts.ListOptions,
			Type:        "all",
		}
		list, resp, err := client.List(context.Background(), org, &prLO)
		logrus.Info("Repos: ", list)
		return list, resp, err
	}
}

func GetPRPaginator(client GitHubPullRequestsClient, org, repo string) func(opts *RateLimitedPage) ([]*github.PullRequest, *github.Response, error) {
	return func(opts *RateLimitedPage) ([]*github.PullRequest, *github.Response, error) {
		prLO := github.PullRequestListOptions{
			ListOptions: opts.ListOptions,
			State:       config.PrState,
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
