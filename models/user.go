package models

type User struct {
	ID       uint32 `json:"id"`
	Username string `json:"login"`
	URL      string `json:"url"`
}
