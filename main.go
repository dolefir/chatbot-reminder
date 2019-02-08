package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/simplewayua/chatbot-reminder/controllers"
	"github.com/simplewayua/chatbot-reminder/db"
	"github.com/simplewayua/chatbot-reminder/models"
	"log"
)

func main() {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	models.AutoMigrate()
	defer db.CloseDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	fmt.Println("Started listening ... 🚀🚀🚀")
	r.POST("/bot", controllers.BotRequestHandler)
	r.Run() // listen and serve on 0.0.0.0:8080
}
