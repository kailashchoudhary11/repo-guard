package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kailashchoudhary11/repo-guard/handlers"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/templates"
)

func main() {
	initializers.LoadDotEnv()
	initializers.LoadGithubClient()
	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		appName := os.Getenv("APP_NAME")
		authorizationUrl := fmt.Sprintf("https://github.com/apps/%v/installations/new", appName)

		template := templates.HomePage(authorizationUrl)
		template.Render(context.Background(), w)
	})

	router.HandleFunc("/webhook", handlers.Webhook)
	fmt.Println("Service is up and running at port 8000")
	http.ListenAndServe(":8000", router)
}
