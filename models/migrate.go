package models

import (
	"github.com/dolefir/chatbot-reminder/db"
)

// AutoMigrate ...
func AutoMigrate() {
	var getDB = db.GetDB()
	getDB.AutoMigrate(&User{})
	getDB.AutoMigrate(&Reminder{})
}
