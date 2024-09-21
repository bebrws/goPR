package store

const (
	PerPage = 10 // Drop this down to test out the pagination!
	PrState = "open"
	ConfigFileName = ".goPR.json"
	StateFileName = ".goPRState.json"
)

// Configuration for the app
type Repo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}
type Config struct {
	Repos []Repo `json:"repos"`
}
