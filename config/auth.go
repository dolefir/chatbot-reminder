package config

import (
	"github.com/simplewayua/chatbot-reminder/models"
	"os"
)

var dp models.DialogFlowProcessor

func AuthDialogFlow() (*models.DialogFlowProcessor, error) {
	dp, err := dp.Init(
		os.Getenv("PROJECT_ID"),
		os.Getenv("AUTH_JSON_FILE_PATH"),
		os.Getenv("LANG"),
		os.Getenv("TIME_ZONE"),
	)
	if err != nil {
		return nil, err
	}
	return dp, nil
}
