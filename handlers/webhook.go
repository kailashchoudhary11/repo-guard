package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside webhook")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	webhookPayload := models.WebhookPayload{}
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		fmt.Println("There was an error in converting json", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Webhook payload", webhookPayload)
	if webhookPayload.Action == "opened" {
		fmt.Println("New issue opened")
		if webhookPayload.Issue.AuthorAssociation == "OWNER" {
			fmt.Println("Issue is opened by repo owner, skipping checks")
			return
		}
		if isSpamIssue := validateIssue(webhookPayload.Repository, &webhookPayload.Issue); isSpamIssue {
			fmt.Println("The issue is spam, closing it")
		}
	}
}

func validateIssue(repo models.Repository, currentIssue *models.Issue) bool {
	fmt.Println("Validating the issue", currentIssue.Number)
	allOpenIssues := services.FetchIssues(initializers.GithubClient, repo)
	fmt.Println("All the open repo issues are ", allOpenIssues)
	for _, issue := range allOpenIssues {
		if issue.Number == currentIssue.Number {
			continue
		}

		if isDuplicate := compareIssues(currentIssue, issue); isDuplicate {
			fmt.Println("The current issue is duplicate")
		}
	}
	return false
}

func compareIssues(issueOne *models.Issue, issueTwo *models.Issue) bool {
	return issueOne.Body == issueTwo.Body
}
