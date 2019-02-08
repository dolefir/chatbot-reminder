package models

import (
	"github.com/simplewayua/chatbot-reminder/db"
)

// AutoMigrate ...
func AutoMigrate() {
	var getDB = db.GetDB()
	getDB.AutoMigrate(&User{})
	getDB.AutoMigrate(&Reminder{})
}
