package main

import (
	"fmt"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/handlers"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service is up and running")
	})
	router.HandleFunc("/webhook", handlers.Webhook)
	fmt.Println("Service is up and running at port 8000")
	http.ListenAndServe(":8000", router)
}
