package main

import (
	"fmt"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/handlers"
	"github.com/kailashchoudhary11/repo-guard/initializers"
)

func main() {
	initializers.LoadDotEnv()
	initializers.LoadGithubClient()

	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/privacy-policy", handlers.PrivacyPolicy)
	router.HandleFunc("/webhook", handlers.Webhook)
	router.HandleFunc("/", handlers.Index)

	fmt.Println("Service is up and running at port 8000")

	http.ListenAndServe(":8000", router)
}
