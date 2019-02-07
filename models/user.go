package models

import (
	"github.com/simplewayua/chatbot-reminder/db"
	"time"
)

// User model
type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"NOT NULL"`
	CreatedAt time.Time
}

// SaveData ...
func (u *User) SaveData() (err error) {
	var getDB = db.GetDB()
	return getDB.Save(u).Error
}
