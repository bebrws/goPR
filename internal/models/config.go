package models

// Configuration for the app
type Repo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}
type Config struct {
	GHToken string `json:"ghtoken"`
	Repos   []Repo `json:"repos"`
}
