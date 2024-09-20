package config

const (
	PerPage = 10 // Drop this down to test out the pagination!
	PrState = "open"
	ConfigFileName = ".goPR.json"
)

type Repo struct {
	Org  string
	Repo string
}
type Config struct {
	Repos []Repo `json:"repos"`
}

type PRReviewComment  struct {
	Login string `json:"commentLogin"`
	Body string `json:"commentBody"`
}

type PRReview  struct {
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
}
type CurrentState struct {
	RepoStates []RepoState `json:"repoStates"`
}

