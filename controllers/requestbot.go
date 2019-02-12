package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mileusna/facebook-messenger"
	"github.com/simplewayua/chatbot-reminder/config"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"github.com/simplewayua/chatbot-reminder/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

// var dp dialogflowmap.DialogFlowProcessor
var r *models.Reminder
var msng *messenger.Messenger
var dp *dialogflowmap.DialogFlowProcessor

// VerificationWebhookHandler ...
func VerificationWebhookHandler(c *gin.Context) {
	challenge := c.Query("hub.challenge")
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")

	if mode == "" && token != os.Getenv("VERIFY_TOKEN") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error, wrong validation token"})
		return
	}

	c.Data(200, "value", []byte(challenge))
}

// SetMessenger ...
func SetMessenger(m *messenger.Messenger) {
	msng = m

	var err error
	dp, err = config.AuthDialogFlow()
	if err != nil {
		panic(err)
	}
}

// MessagesWebhookHandler ...
func MessagesWebhookHandler(c *gin.Context) {
	msng.VerifyWebhook(c.Writer, c.Request)            // verify webhook if asked from Facebook
	fbRequest, _ := messenger.DecodeRequest(c.Request) // decode entire request received from Facebook into FacebookRequest struct

	for _, entry := range fbRequest.Entry {
		// pageID := entry.ID  // here you can find page id that received message
		for _, msg := range entry.Messaging {
			userID := msg.Sender.ID // user that sent you a message
			user := &models.User{}

			switch {
			case msg.Message != nil:
				// log.Println("Msg received with content:", msg.Message.Text)
				if len(msg.Message.Text) == 0 {
					return
				}
				// Use NPLResponse
				res, _ := dp.ProcessNPL(msg.Message.Text, strconv.Itoa(int(userID)))
				if res.Intent == "GetDateWithTime" {
					result, err := models.DialoglowResponseToReminder(&res)
					if err != nil {
						msng.SendTextMessage(userID, "Sorry have a trouble! :(")
					}

					result.NameID = strconv.Itoa(int(userID))
					result.SaveData()

					user.Name = strconv.Itoa(int(userID))
					user.SaveData()
				}

				if res.Intent == "GetListReminder" {
					for _, k := range res.Entities {
						lists, err := models.GetTimesReminder(k)
						if err != nil {
							continue
						}

						if len(lists) == 0 {
							msng.SendTextMessage(userID, "You don't have reminders for the selected time.")
							return
						}

						for _, text := range lists {
							msng.SendTextMessage(userID, text.Text)
						}
					}
				}

				if res.Intent == "DeleteReminder" {
					for _, k := range res.Entities {
						models.DeleteReminder(k)
					}
				}

				msng.SendTextMessage(userID, res.FulfillmentText)

			case msg.Delivery != nil:
				// delivery report received, check First example what to do next

			case msg.Postback != nil:
				// postback received, check First example what can you do with that
				log.Println("Postback received with content:", msg.Postback.Payload)
				if msg.Postback.Payload == "ACCEPT" {
					msng.SendTextMessage(userID, "Thank")
				}
				if msg.Postback.Payload == "SNOOZE" {
					_, err := models.UpdateReminderNotificationTime()
					if err != nil {
						panic(err)
					}
					msng.SendTextMessage(userID, "Ok, I will postpone the reminder for 10 minutes")
				}
			}
		}
	}
}
