package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/facebook-messenger"
	"github.com/simplewayua/chatbot-reminder/config"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"github.com/simplewayua/chatbot-reminder/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var dp dialogflowmap.DialogFlowProcessor
var r models.Reminder

// BotRequestHandler for testing postman
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

func MessagesWebhookHandler(c *gin.Context) {
	msng := &messenger.Messenger{
		AccessToken: os.Getenv("PAGE_ACCESS_TOKEN"),
		VerifyToken: os.Getenv("VERIFY_TOKEN"),
		PageID:      os.Getenv("PAGE_ID"),
	}

	dp, err := config.AuthDialogFlow()
	if err != nil {
		panic(err)
	}

	msng.VerifyWebhook(c.Writer, c.Request)            // verify webhook if asked from Facebook
	fbRequest, _ := messenger.DecodeRequest(c.Request) // decode entire request received from Facebook into FacebookRequest struct

	// now you have it all and you can do whatever you want with received request
	// enumerate each entry, and each message in entry
	for _, entry := range fbRequest.Entry {
		// pageID := entry.ID  // here you can find page id that received message
		for _, msg := range entry.Messaging {
			userID := msg.Sender.ID // user that sent you a message
			user := &models.User{}
			// but "message" can be text message, delivery report or postback, so check it what it is
			// it can only be one of this, so we use switch
			switch {
			case msg.Message != nil:
				log.Println("Msg received with content:", msg.Message.Text)
				// msng.SendTextMessage(userID, "Hello there")
				// check First example for more sending messages examples

				// Use NPLResponse
				res := dp.ProcessNPL(msg.Message.Text, strconv.Itoa(int(userID)))
				if res.Intent == "GetDateWithTime" {
					result, err := r.ToReminder(&res)
					if err != nil {
						msng.SendTextMessage(userID, "Sorry have a trouble! :(")
					}
					result.NameID = strconv.Itoa(int(userID))
					result.SaveData()
					user.Name = strconv.Itoa(int(userID))
					// user.Reminder = *result
					user.SaveData()
				}

				if res.Intent == "GetListReminder" {
					for _, k := range res.Entities {
						lists := r.GetTimesReminder(k)
						for _, text := range lists {
							msng.SendTextMessage(userID, text.Text)
						}
					}
				}
				if res.Intent == "DeleteReminder" {
					for _, k := range res.Entities {
						r.DeleteReminder(k)
					}
				}
				msng.SendTextMessage(userID, res.FulfillmentText)

			case msg.Delivery != nil:
				// delivery report received, check First example what to do next

			case msg.Postback != nil:
				// postback received, check First example what can you do with that
				log.Println("Postback received with content:", msg.Postback.Payload)
			}
		}
	}
}
