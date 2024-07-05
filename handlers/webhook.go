package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside webhook")
	clientId := os.Getenv("CLIENT_ID")
	jwtToken, err := services.GenerateJWTForApp(clientId, "repository-guard.2024-07-02.private-key.pem")
	fmt.Println("The token is", jwtToken)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
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

	accessToken := services.GetInstallationAccessToken(webhookPayload.Installation.ID, jwtToken)
	authenticatedClient := initializers.GetClientWithToken(accessToken)

	if webhookPayload.Action == "opened" {
		fmt.Println("New issue opened")
		// if webhookPayload.Issue.AuthorAssociation == "OWNER" {
		// 	fmt.Println("Issue is opened by repo owner, skipping checks")
		// 	return
		// }
		if isSpamIssue := validateIssue(authenticatedClient, webhookPayload.Repository, &webhookPayload.Issue); isSpamIssue {
			fmt.Println("The duplicate issue is closed successfully")
		}
	}
}

func validateIssue(githubClient *github.Client, repo models.Repository, currentIssue *models.Issue) bool {
	fmt.Println("Validating the issue", currentIssue.Number)

	duplicateIssue := make(chan int)

	allOpenIssues := services.FetchIssues(githubClient, repo)
	for _, issue := range allOpenIssues {
		if issue.Number == currentIssue.Number {
			continue
		}

		go compareIssues(currentIssue, issue, duplicateIssue)

	}
	for i := 0; i < len(allOpenIssues); i++ {
		issueNumber := <-duplicateIssue
		if issueNumber > -1 {
			fmt.Printf("The issue is duplicate of %v, closing the issue.\n", issueNumber)
			closingReason := fmt.Sprintf("A similar issue already exists. Please check #%v", issueNumber)
			err := services.CloseIssue(githubClient, repo, currentIssue.Number, closingReason)
			if err != nil {
				fmt.Println("Error in closing the issue", err)
				return false
			}
			return true
		}
	}
	return false
}

func compareIssues(issueOne *models.Issue, issueTwo *models.Issue, isDuplicate chan int) {
	fmt.Println("Comparing the issues")
	payload := fmt.Sprintf(`{"issue1_title": "%v", "issue1_body": "", "issue2_title": "%v", "issue2_body": "" }`, issueOne.Title, issueTwo.Title)
	jsonBody := []byte(payload)

	bodyReader := bytes.NewReader(jsonBody)
	response := struct {
		Similarity float32 `json:"similarity"`
	}{}

	AIModelURL := os.Getenv("AI_MODEL_URL")

	requestURL := fmt.Sprintf("%vcompare_issues", AIModelURL)
	res, err := http.Post(requestURL, "application/json", bodyReader)
	fmt.Println("Res", res)
	if err != nil {
		fmt.Println("Error in making compare issues request", err)
		isDuplicate <- -1
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Cannot read response body", err)
		isDuplicate <- -1
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Response is in invalid format", err)
		isDuplicate <- -1
	}
	fmt.Println("Response", res)
	if response.Similarity > 0.75 {
		isDuplicate <- issueTwo.Number
	}

	isDuplicate <- -1
}
