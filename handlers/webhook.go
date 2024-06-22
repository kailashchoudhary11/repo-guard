package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/models"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	jsonBody := models.WebhookPayload{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		fmt.Println("There was an error in converting json", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("issue details are ", jsonBody)
}
