package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/simplewayua/chatbot-reminder/config"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"io/ioutil"
	"net/http"
)

var dp dialogflowmap.DialogFlowProcessor

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

	// Use NLP
	response := dp.ProcessNPL(m.Message, "UserName")
	c.JSON(http.StatusOK, response)
}
