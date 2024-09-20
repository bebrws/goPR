package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/notify"
	"github.com/google/go-github/v65/github"
)

type PaginateAble interface {
	github.ListOptions
}

// What is a RepositoryComment??
// https://github.com/google/go-github/blob/23592519310b534946cbb8527ad6a92d5d29ddd0/github/event_types.go#L71C6-L71C24
// github.RepositoryComment Is inside CommitCommentEvent 
// which `CommitCommentEvent is triggered when a commit comment is created.`

type PaginateAbleGithubTypes interface {
	*github.PullRequest | *github.Comment | *github.PullRequestComment | *github.PullRequestReview | *github.RepositoryComment
}

func paginate[R PaginateAbleGithubTypes](pf func() ([]R, error)) ([]R, error) {
	allItems := []R{}
	for {
		items, err := pf()
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

func getPRPaginator(client *github.Client, org, repo string) func() ([]*github.PullRequest, error) {
	prReqOpts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: config.PerPage},
		State: 	 config.PrState,
	}
	return func() ([]*github.PullRequest, error) {
		// Get the PRs
		fmt.Println("Making page request wiht: ", prReqOpts.Page)
		prs, resp, err := client.PullRequests.List(context.Background(), org, repo, prReqOpts)
		if err != nil {
			fmt.Printf("Error listing Review Comments for %s/%s: %s", org, repo, err)
			return nil, err
		}
		if resp != nil {
			fmt.Println("Rate limit remaining:", resp.Rate.Remaining)
		}
		if len(prs) == 0 {
			fmt.Println("No more PRs found")
			return prs, nil
		}
		prReqOpts.Page = resp.NextPage
		return prs, nil
	}
}

func getReviewPaginator(client *github.Client, org, repo string, prNumber int) func() ([]*github.PullRequestReview, error) {
	listOps := github.ListOptions{PerPage: config.PerPage}
	return func() ([]*github.PullRequestReview, error) {
		revs, resp, err := client.PullRequests.ListReviews(context.Background(), org, repo, prNumber, &listOps)
		if err != nil {
			fmt.Printf("Error listing Review Comments for %s/%s: %s", org, repo, err)
			return revs, err
		}
		if resp != nil {
			fmt.Println("Rate limit remaining:", resp.Rate.Remaining)
		}
		if len(revs) == 0 {
			fmt.Println("No more PRs found")
			return revs, nil
		}

		listOps.Page = resp.NextPage
		return revs, nil
	}
}

func getReviewCommentsPaginator(client *github.Client, org, repo string, prNumber int, revID int64) func() ([]*github.PullRequestComment, error) {
	listOps := github.ListOptions{PerPage: config.PerPage}
	return func() ([]*github.PullRequestComment, error) {
		revComments, resp, err := client.PullRequests.ListReviewComments(context.Background(), config.Org, config.Repo, prNumber, revID, &listOps)
		if err != nil {
			fmt.Printf("Error listing Review Comments for %s/%s: %s", org, repo, err)
			return revComments, err
		}
		if resp != nil {
			fmt.Println("Rate limit remaining:", resp.Rate.Remaining)
		}
		if len(revComments) == 0 {
			fmt.Println("No Review Comments found")
			return revComments, nil
		}

		listOps.Page = resp.NextPage
		return revComments, nil
	}
}

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

	allPrs, err := paginate(getPRPaginator(client, config.Org, config.Repo))
	if err != nil {
		fmt.Println("Error getting PRs", err)
		os.Exit(1)
	}

	for i, pr := range allPrs {
		fmt.Printf("%d. PR title: %s\n", i+1, *pr.Title)

		allRevs, err := paginate(getReviewPaginator(client, config.Org, config.Repo, *pr.Number))
		if err != nil {
			fmt.Println("Error getting PR Reviews", err)
			os.Exit(1)
		}

		fmt.Println("Reviews for PR: ", *pr.Title, " ---- # ", *pr.Number)
		for j, rev := range allRevs {
			fmt.Printf("  %d. Review: %s from %s\n", j+1, *rev.Body, rev.User.GetLogin())

			allRevs, err := paginate(getReviewCommentsPaginator(client, config.Org, config.Repo, *pr.Number, *rev.ID))
			if err != nil {
				fmt.Println("Error getting PR Reviews", err)
				os.Exit(1)
			}

			for j, rev := range allRevs {
				fmt.Printf("      %d. Review Comment: %s from %s\n", j+1, *rev.Body, rev.User.GetLogin())
			}
		}
	}

}
