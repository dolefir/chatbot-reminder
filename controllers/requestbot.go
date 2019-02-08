package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/simplewayua/chatbot-reminder/config"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"github.com/simplewayua/chatbot-reminder/models"
	"io/ioutil"
	"net/http"
)

var dp dialogflowmap.DialogFlowProcessor
var r models.Reminder

// BotRequestHandler ...
func BotRequestHandler(c *gin.Context) {
	dp, err := config.AuthDialogFlow()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error reading request body"})
	}

	type inboundMessage struct {
		Message string
	}

	var m inboundMessage
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}

	// Use NPLResponse
	NPLResponse := dp.ProcessNPL(m.Message, "UserName")
	if NPLResponse.Intent == "GetDateWithTime" {
		result, err := r.ToReminder(&NPLResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error reading request body"})
		}
		result.SaveData()
	}

	c.JSON(http.StatusOK, NPLResponse)
}
