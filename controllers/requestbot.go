package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/simplewayua/chatbot-reminder/config"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"github.com/simplewayua/chatbot-reminder/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// VerificationWebhookHandler ...
func VerificationWebhookHandler(c *gin.Context) {
	challenge := c.Query("hub.challenge")
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")

	if mode != "" && token == os.Getenv("VERIFY_TOKEN") {
		c.Data(200, "value", []byte(challenge))
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error, wrong validation token"})
	}
}

// MessagesWebhookHandler ...
func MessagesWebhookHandler(c *gin.Context) {
	var callback models.Callback
	json.NewDecoder(c.Request.Body).Decode(&callback)
	log.Println(callback.Object)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				log.Println(event.Message.Text)
			}
		}
		c.Data(http.StatusOK, "message", []byte("Got your message"))
	} else {
		c.Data(http.StatusNotFound, "message", []byte("Message not supported"))
	}
}
