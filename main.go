package main

import (
	"fmt"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/handlers"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
)

func main() {
	initializers.LoadDotEnv()
	initializers.LoadGithubClient()
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service is up and running")
	})
	router.HandleFunc("/webhook", handlers.Webhook)
	router.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		repo := models.Repository{
			Name: "Moksh",
			Owner: models.User{
				Username: "akshitagupta15june",
			},
		}
		issues := services.FetchIssues(initializers.GithubClient, repo)
		fmt.Println("The number of issues are ", len(issues))
	})
	fmt.Println("Service is up and running at port 8000")
	http.ListenAndServe(":8000", router)
}
