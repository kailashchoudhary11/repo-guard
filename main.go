package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service is up and running")
	})
	http.ListenAndServe(":8000", router)
}
