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

type batchIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type batchRequest struct {
	CurrentIssue batchIssue   `json:"current_issue"`
	OtherIssues  []batchIssue `json:"other_issues"`
	Threshold    float32      `json:"threshold"`
}

type similarIssueResult struct {
	Number     int     `json:"number"`
	Similarity float32 `json:"similarity"`
}

type batchResponse struct {
	SimilarIssues []similarIssueResult `json:"similar_issues"`
}

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
		config := &models.InstallationConfig{
			ShouldClose: false,
			Language:    "english",
			Sensitivity: 90,
		}

		conn, err := initializers.GetDBClient(ctx)
		if err != nil {
			fmt.Println("Unable to connect to database:", err)
		} else {
			defer conn.Close(ctx)
			queries := db.New(conn)
			installationDetails, err := queries.GetInstallationByInstallationID(ctx, fmt.Sprintf("%d", webhookPayload.Installation.ID))
			if err != nil {
				fmt.Println("Error in fetching installation details", err)
			} else {
				dbConfig, err := installationDetails.Config()
				if err != nil {
					fmt.Println("Error in parsing installation config", err)
				} else {
					config = dbConfig
				}
			}
		}

		threshold := float32(config.Sensitivity) / 100.0
		handleSimilarityCheck(context.Background(), authenticatedClient, webhookPayload.Repository, &webhookPayload.Issue, config.ShouldClose, threshold)
	}
}

func handleSimilarityCheck(
	ctx context.Context,
	githubClient *github.Client,
	repo models.Repository,
	currentIssue *models.Issue,
	shouldClose bool,
	threshold float32,
) {
	allOpenIssues := services.FetchIssues(githubClient, repo)

	otherIssues := make([]batchIssue, 0, len(allOpenIssues))
	for _, issue := range allOpenIssues {
		if issue.Number == currentIssue.Number {
			continue
		}
		otherIssues = append(otherIssues, batchIssue{
			Number: issue.Number,
			Title:  issue.Title,
			Body:   issue.Body,
		})
	}

	if len(otherIssues) == 0 {
		fmt.Println("No other open issues to compare against")
		return
	}

	reqPayload := batchRequest{
		CurrentIssue: batchIssue{
			Number: currentIssue.Number,
			Title:  currentIssue.Title,
			Body:   currentIssue.Body,
		},
		OtherIssues: otherIssues,
		Threshold:   threshold,
	}

	jsonBody, err := json.Marshal(reqPayload)
	if err != nil {
		fmt.Println("Error marshaling batch request", err)
		return
	}

	aiModelURL := os.Getenv("AI_MODEL_URL")
	requestURL := fmt.Sprintf("%vbatch_compare", aiModelURL)

	res, err := http.Post(requestURL, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("Error in batch compare request", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Cannot read batch response body", err)
		return
	}

	var batchResp batchResponse
	if err := json.Unmarshal(body, &batchResp); err != nil {
		fmt.Println("Batch response is in invalid format", err)
		return
	}

	if len(batchResp.SimilarIssues) == 0 {
		fmt.Println("No similar issues found")
		return
	}

	issueLinks := make([]string, 0, len(batchResp.SimilarIssues))
	for _, si := range batchResp.SimilarIssues {
		issueLinks = append(issueLinks, fmt.Sprintf("#%d", si.Number))
	}

	issueComment := fmt.Sprintf(
		"⚠️ It looks like similar open issues already exist: %s.\n"+
			"Please check these before proceeding to avoid duplicates.",
		strings.Join(issueLinks, ", "),
	)

	fmt.Println(issueComment)

	if err := services.AddComment(ctx, githubClient, repo, currentIssue.Number, issueComment); err != nil {
		fmt.Println("Error in adding comment on the issue", err)
		return
	}

	if shouldClose {
		if err := services.CloseIssue(ctx, githubClient, repo, currentIssue.Number); err != nil {
			fmt.Println("Error in closing the issue", err)
			return
		}
	}
}
