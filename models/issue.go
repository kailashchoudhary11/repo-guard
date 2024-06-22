package models

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
