package db

import (
	"encoding/json"

	"github.com/kailashchoudhary11/repo-guard/models"
)

func (i Installation) Config() (*models.InstallationConfig, error) {
	var cfg models.InstallationConfig
	if err := json.Unmarshal(i.ConfigData, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
