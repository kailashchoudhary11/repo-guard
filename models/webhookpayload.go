package models

type WebhookPayload struct {
	Action     string     `json:"string"`
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
}
