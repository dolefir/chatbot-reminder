package models

import (
	"time"

	"github.com/dolefir/chatbot-reminder/db"
)

// User has many Reminders, ID is the foreign key
type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"unique;not null"`
	CreatedAt time.Time
}

// SaveData ...
func (u *User) SaveData() (err error) {
	var getDB = db.GetDB()
	return getDB.Save(u).Error
}
