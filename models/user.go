package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"login"`
	URL      string `json:"url"`
}
