package config

import "time"

const (
	PerPage = 10 // Drop this down to test out the pagination!
	PrState = "open"
	ConfigFileName = ".goPR.json"
	StateFileName = ".goPR.json"
)

type Repo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}
type Config struct {
	Repos []Repo `json:"repos"`
}

type PRReviewComment  struct {
	ID int64 `json:"commentID"`
	UpdatedAt time.Time `json:"commentUpdatedAt"`
	Login string `json:"commentLogin"`
	Body string `json:"commentBody"`
}

type PRReview  struct {
	ID int64 `json:"reviewID"`
	Login string `json:"reviewLogin"`
	Body string `json:"reviewBody"`
	Comments []PRReviewComment `json:"reviewComments"`
}
type PR struct {
	Number int `json:"prNumber"`
	Body string `json:"prBody"`
	Reviews []PRReview `json:"prReviews"`
}

type RepoState struct {
	PRs []PR `json:"prs"`
	Name string `json:"repoName"`
}
type GHState struct {
	RepoStates []RepoState `json:"repoStates"`
}

