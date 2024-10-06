package models

import (
	"time"
)

// State for the app
type PRReviewComment struct {
	ID        int64     `json:"commentID"`
	UpdatedAt time.Time `json:"commentUpdatedAt"`
	Login     string    `json:"commentLogin"`
	Body      string    `json:"commentBody"`
}

type PRReview struct {
	ID       int64             `json:"reviewID"`
	Login    string            `json:"reviewLogin"`
	Body     string            `json:"reviewBody"`
	Comments []PRReviewComment `json:"reviewComments"`
}
type PR struct {
	Number  int        `json:"prNumber"`
	Body    string     `json:"prBody"`
	Reviews []PRReview `json:"prReviews"`
}

type RepoState struct {
	PRs  []PR   `json:"prs"`
	Name string `json:"repoName"`
}
type GHState struct {
	RepoStates []RepoState `json:"repoStates"`
}
