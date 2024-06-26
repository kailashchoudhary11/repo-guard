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

func _fetchIssuesForPage(client *github.Client, repo models.Repository, page int, ch chan []*models.Issue, countCh chan int, sendCount bool) {
	convertedIssues := []*models.Issue{}
	opts := &github.IssueListByRepoOptions{State: "open", ListOptions: github.ListOptions{Page: page, PerPage: 100}}
	issues, res, err := client.Issues.ListByRepo(context.Background(), repo.Owner.Username, repo.Name, opts)
	if err != nil {
		fmt.Println("The error is", err)
		return
	}
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}
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
	ch <- convertedIssues
	if sendCount {
		countCh <- res.LastPage
	}
}

func FetchIssues(client *github.Client, repo models.Repository) []*models.Issue {
	fmt.Println("Fetching the issues")
	allIssues := []*models.Issue{}
	ch := make(chan []*models.Issue)
	countCh := make(chan int)
	go _fetchIssuesForPage(client, repo, 1, ch, countCh, true)
	allIssues = append(allIssues, <-ch...)
	lastPage := <-countCh
	if lastPage == 0 {
		return allIssues
	}
	for i := 2; i <= lastPage; i++ {
		go _fetchIssuesForPage(client, repo, i, ch, countCh, false)
	}
	for i := 2; i <= lastPage; i++ {
		allIssues = append(allIssues, <-ch...)
	}
	fmt.Println("Fetched the issues")
	return allIssues
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
