package models

import (
	"github.com/simplewayua/chatbot-reminder/db"
	"time"
)

// User has many Reminders, ID is the foreign key
type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"NOT NULL"`
	CreatedAt time.Time
	Reminders []Reminder `gorm:"foreignkey:ID"`
}

// SaveData ...
func (u *User) SaveData() (err error) {
	var getDB = db.GetDB()
	return getDB.Save(u).Error
}
