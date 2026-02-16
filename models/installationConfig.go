package models

type InstallationConfig struct {
	ShouldClose bool   `json:"should_close"`
	Language    string `json:"language"`
	Sensitivity int `json:"sensitivity"`
}
