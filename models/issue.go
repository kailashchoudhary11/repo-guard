package models

type Issue struct {
	ID                int64  `json:"id"`
	Number            int    `json:"number"`
	State             string `json:"state"`
	Title             string `json:"title"`
	Body              string `json:"body"`
	URL               string `json:"url"`
	RepositoryURL     string `json:"repository_url"`
	LabelsURL         string `json:"labels_url"`
	Author            User   `json:"user"`
	AuthorAssociation string `json:"author_association"`
}
