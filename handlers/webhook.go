package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	db "github.com/kailashchoudhary11/repo-guard/db/generated"
	"github.com/kailashchoudhary11/repo-guard/helpers"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("CLIENT_ID")
	privatePem := os.Getenv("PRIVATE_KEY")
	privatePem = strings.ReplaceAll(privatePem, "\\n", "\n")
	jwtToken, err := helpers.GenerateJWT(clientId, privatePem)
	ctx := r.Context()
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

	if webhookPayload.Action == "opened" {
		accessToken := services.GetInstallationAccessToken(webhookPayload.Installation.ID, jwtToken)
		authenticatedClient := initializers.GetClientWithToken(accessToken)
		shouldClose := false

		conn, err := initializers.GetDBClient(ctx)
		if err != nil {
			fmt.Println("Unable to connect to database:", err)
		}
		defer conn.Close(ctx)
		queries := db.New(conn)
		installationDetails, err := queries.GetInstallationByInstallationID(ctx, string(webhookPayload.Installation.ID))
		config, err := installationDetails.Config()
		if err != nil {
			fmt.Println("Error in fetching installation config", err)
		} else {
			shouldClose = config.ShouldClose
		}

		handleSimilarityCheck(r.Context(), authenticatedClient, webhookPayload.Repository, &webhookPayload.Issue, shouldClose)
	}
}

func handleSimilarityCheck(ctx context.Context, githubClient *github.Client, repo models.Repository, currentIssue *models.Issue, shouldClose bool) bool {
	duplicateIssue := make(chan int)

	allOpenIssues := services.FetchIssues(githubClient, repo)
	for _, issue := range allOpenIssues {
		if issue.Number == currentIssue.Number {
			continue
		}

		go compareIssues(currentIssue, issue, duplicateIssue)

	}
	similarIssues := []int{}
	issueComment := fmt.Sprintf("Similar open issues already exist. Please check")
	for i := 0; i < len(allOpenIssues); i++ {
		issueNumber := <-duplicateIssue
		if issueNumber > -1 {
			issueComment = fmt.Sprintf("%v #%v,", issueComment, issueNumber)
			similarIssues = append(similarIssues, issueNumber)
		}
	}
	if len(similarIssues) == 0 {
		fmt.Println("No similar issues found")
		return false
	}

	fmt.Printf("The similar issues are", similarIssues)

	if err := services.AddComment(ctx, githubClient, repo, currentIssue.Number, issueComment); err != nil {
		fmt.Println("Error in adding comment on the issue", err)
		return false
	}
	if shouldClose {
		if err := services.CloseIssue(ctx, githubClient, repo, currentIssue.Number); err != nil {
			fmt.Println("Error in closing the issue", err)
			return false
		}
	}
	return true
}

func compareIssues(issueOne *models.Issue, issueTwo *models.Issue, isDuplicate chan int) {
	fmt.Println("Will be now comparing the issues")
	payload := fmt.Sprintf(`{"issue1_title": "%v", "issue1_body": "", "issue2_title": "%v", "issue2_body": "" }`, issueOne.Title, issueTwo.Title)
	jsonBody := []byte(payload)

	bodyReader := bytes.NewReader(jsonBody)
	response := struct {
		Similarity float32 `json:"similarity"`
	}{}

	AIModelURL := os.Getenv("AI_MODEL_URL")

	requestURL := fmt.Sprintf("%vcompare_issues", AIModelURL)
	res, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		fmt.Println("Error in making compare issues request", err)
		isDuplicate <- -1
	}

	body, err := io.ReadAll(res.Body)
	fmt.Println("The response from API is", string(body))
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Cannot read response body", err)
		isDuplicate <- -1
	}
	fmt.Println(response)

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Response is in invalid format", err)
		isDuplicate <- -1
	}
	if response.Similarity > 0.75 {
		isDuplicate <- issueTwo.Number
	}

	isDuplicate <- -1
}
