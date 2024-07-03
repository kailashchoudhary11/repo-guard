package initializers

import (
	"os"

	"github.com/google/go-github/v62/github"
)

var GithubClient *github.Client

func LoadGithubClient() {
	token, isTokenPresent := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if isTokenPresent {
		GithubClient = github.NewClient(nil).WithAuthToken(token)
	} else {
		GithubClient = github.NewClient(nil)
	}
}

func GetClientWithToken(accessToken string) *github.Client {
	return github.NewClient(nil).WithAuthToken(accessToken)
}
