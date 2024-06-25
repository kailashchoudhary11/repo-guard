package services

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
	"github.com/kailashchoudhary11/repo-guard/models"
)

type ProjectDetail struct {
	owner    string
	repoName string
}

func FetchIssues(client *github.Client, repo models.Repository) []*models.Issue {
	fmt.Println("Fetching issues")
	opts := &github.IssueListByRepoOptions{State: "open", ListOptions: github.ListOptions{PerPage: 100}}
	issues, _, err := client.Issues.ListByRepo(context.Background(), repo.Owner.Username, repo.Name, opts)
	if err != nil {
		fmt.Println("The error is", err)
		return nil
	}
	fmt.Println("Fetched the issues")
	convertedIssues := []*models.Issue{}
	for _, issue := range issues {
		user := models.User{
			Username: *issue.User.Login,
			ID:       *issue.User.ID,
			URL:      *issue.User.URL,
		}
		convertedIssue := models.Issue{
			ID:                *issue.ID,
			Number:            *issue.Number,
			State:             *issue.State,
			Title:             *issue.Title,
			URL:               *issue.URL,
			RepositoryURL:     *issue.RepositoryURL,
			LabelsURL:         *issue.LabelsURL,
			Author:            user,
			AuthorAssociation: *issue.AuthorAssociation,
		}
		if issue.Body != nil {
			convertedIssue.Body = *issue.Body
		}
		convertedIssues = append(convertedIssues, &convertedIssue)
	}
	return convertedIssues
}

func CloseIssue(client *github.Client, repo models.Repository, issueNumber int, reason string) error {
	if reason != "" {
		issueComment := github.IssueComment{
			Body: &reason,
		}
		_, _, err := client.Issues.CreateComment(context.Background(), repo.Owner.Username, repo.Name, issueNumber, &issueComment)
		if err != nil {
			fmt.Println("Error in adding comment", err)
			return err
		}
	}

	state := "closed"
	issueRequest := github.IssueRequest{
		State: &state,
	}
	_, _, err := client.Issues.Edit(context.Background(), repo.Owner.Username, repo.Name, issueNumber, &issueRequest)
	if err != nil {
		fmt.Println("Error in closing issue", err)
		return err
	}
	return nil
}
