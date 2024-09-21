package store

// Configuration for the app
type Repo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}
type Config struct {
	Repos []Repo `json:"repos"`
}
