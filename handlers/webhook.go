package handlers

import (
	"bytes"
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
	for _, issue := range allOpenIssues {
		if issue.Number == currentIssue.Number {
			continue
		}

		if isDuplicate := compareIssues(currentIssue, issue); isDuplicate {
			fmt.Printf("The issue %v is duplicate\n", issue.Number)
			services.CloseIssue(initializers.GithubClient, repo, currentIssue.Number)
			return true
		}
	}
	return false
}

func compareIssues(issueOne *models.Issue, issueTwo *models.Issue) bool {
	payload := fmt.Sprintf(`{"issue1_title": "%v", "issue1_body": "", "issue2_title": "%v", "issue2_body": "" }`, issueOne.Title, issueTwo.Title)
	jsonBody := []byte(payload)

	bodyReader := bytes.NewReader(jsonBody)
	response := struct {
		Similarity float32 `json:"similarity"`
	}{}

	requestURL := "http://localhost:5000/compare_issues"
	res, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		fmt.Println("Error in making compare issues request", err)
		return false
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Cannot read response body", err)
		return false
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Response is in invalid format", err)
		return false
	}

	return response.Similarity > 0.75
}
