package handlers

import (
	"context"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/templates"
)

func PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	template := templates.PrivacyPolicy()
	template.Render(context.Background(), w)
}
