package models

type WebhookPayload struct {
	Action       string       `json:"action"`
	Issue        Issue        `json:"issue"`
	Repository   Repository   `json:"repository"`
	Installation Installation `json:"installation"`
}
