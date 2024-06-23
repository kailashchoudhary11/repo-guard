package models

type Repository struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    User   `json:"owner"`
}
