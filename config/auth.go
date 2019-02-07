package config

import (
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"os"
)

var dp dialogflowmap.DialogFlowProcessor

func AuthDialogFlow() (*dialogflowmap.DialogFlowProcessor, error) {
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
