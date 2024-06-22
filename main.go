package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

type User struct {
	ID       uint32 `json:"id"`
	Username string `json:"login"`
	URL      string `json:"url"`
}

type Issue struct {
	ID                uint32 `json:"id"`
	Number            uint16 `json:"number"`
	Title             string `json:"title"`
	State             string `json:"state"`
	URL               string `json:"url"`
	RepositoryURL     string `json:"repository_url"`
	LabelsURL         string `json:"labels_url"`
	Author            User   `json:"user"`
	AuthorAssociation string `json:"author_association"`
}

type WebhookPayload struct {
	Action     string     `json:"string"`
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service is up and running")
	})
	router.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		jsonBody := WebhookPayload{}
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			fmt.Println("There was an error in converting json", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("issue details are ", jsonBody)
	})
	fmt.Println("Service is up and running at port 8000")
	http.ListenAndServe(":8000", router)
}
