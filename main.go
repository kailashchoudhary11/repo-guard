package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/handlers"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
	"github.com/kailashchoudhary11/repo-guard/templates"
)

func main() {
	initializers.LoadDotEnv()
	initializers.LoadGithubClient()
	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // if using net/http
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := templates.HomePage()
		index.Render(context.Background(), w)
	})

	router.HandleFunc("/webhook", handlers.Webhook)
	router.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		repo := models.Repository{
			Name: "Flipkart_Clone",
			Owner: models.User{
				Username: "arghadipmanna101",
			},
		}
		issues := services.FetchIssues(initializers.GithubClient, repo)
		fmt.Println("The number of issues are ", len(issues))
	})
	fmt.Println("Service is up and running at port 8000")
	http.ListenAndServe(":8000", router)
}
