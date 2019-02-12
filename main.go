package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mileusna/facebook-messenger"
	"github.com/simplewayua/chatbot-reminder/controllers"
	"github.com/simplewayua/chatbot-reminder/db"
	"github.com/simplewayua/chatbot-reminder/models"
	"github.com/simplewayua/chatbot-reminder/services"
	"log"
	"os"
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
