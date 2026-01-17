package initializers

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var Oauth2Config *oauth2.Config

func LoadOauth2Config() {
	Oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
	}
}
