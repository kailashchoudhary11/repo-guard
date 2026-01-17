package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kailashchoudhary11/repo-guard/templates"
)

func Index(w http.ResponseWriter, r *http.Request) {
	appName := os.Getenv("APP_NAME")
	authorizationUrl := fmt.Sprintf("https://github.com/apps/%v/installations/new", appName)

	template := templates.HomePage(authorizationUrl)
	template.Render(context.Background(), w)
}
