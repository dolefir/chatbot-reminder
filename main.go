package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dolefir/chatbot-reminder/controllers"
	"github.com/dolefir/chatbot-reminder/db"
	"github.com/dolefir/chatbot-reminder/models"
	"github.com/dolefir/chatbot-reminder/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	messenger "github.com/mileusna/facebook-messenger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	msng := &messenger.Messenger{
		AccessToken: os.Getenv("PAGE_ACCESS_TOKEN"),
		VerifyToken: os.Getenv("VERIFY_TOKEN"),
		PageID:      os.Getenv("PAGE_ID"),
	}

	controllers.SetMessenger(msng)

	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	models.AutoMigrate()
	defer db.CloseDB()

	go services.MonitorBotNotification(msng)

	r := gin.Default()
	fmt.Println("Started listening ... 🚀🚀🚀")
	r.GET("/webhook", controllers.VerificationWebhookHandler)
	r.POST("/webhook", controllers.MessagesWebhookHandler)
	r.Run(":3000")
}
