package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/dolefir/chatbot-reminder/db"
	"github.com/dolefir/chatbot-reminder/dialogflowmap"
)

// Reminder ...
type Reminder struct {
	ID       uint   `gorm:"primary_key"`
	Text     string `gorm:"NOT NULL"`
	Time     string
	User     User `gorm:"association_foreignkey:Name"`
	NameID   string
	Position uint
}

// NPL ...
type NPL struct {
	Intent     string   `json:"intent"`
	Confidence float32  `json:"confidence"`
	Entities   Entities `json:"entities"`
}

// Entities ...
type Entities struct {
	Datewithtime         string `json:"datewithtime"`
	DatewithtimeOriginal string `json:"-"`
	Text                 string `json:"text"`
	TextOriginal         string `json:"-"`
}

// SaveData ...
func (r *Reminder) SaveData() (err error) {
	var getDB = db.GetDB()
	return getDB.Save(r).Error
}

// ToReminder ...
func DialoglowResponseToReminder(response *dialogflowmap.NPLResponse) (*Reminder, error) {
	b, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	npl := NPL{}
	err = json.Unmarshal([]byte(b), &npl)
	if err != nil {
		return nil, err
	}

	reminder := Reminder{
		Text:     npl.Entities.Text,
		Time:     npl.Entities.Datewithtime,
		Position: 0,
	}
	return &reminder, nil
}

// GetTimesReminder ...
func GetTimesReminder(t string) (rems []Reminder, err error) {
	var getDB = db.GetDB()
	s := strings.Split(t, "T")
	tm := s[0]
	if err = getDB.Where("time LIKE ? AND position = ?", tm+"%", 0).Find(&rems).Error; err != nil {
		return nil, err
	}
	return
}

// DeleteReminder ...
func DeleteReminder(t string) (err error) {
	var getDB = db.GetDB()
	if err = getDB.Where("text LIKE ?", "%"+t+"%").Delete(Reminder{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllTimeReminder for notification
func GetAllTimeReminder() (rems []Reminder, err error) {
	var getDB = db.GetDB()
	if err := getDB.Find(&rems).Error; err != nil {
		return nil, err
	}
	return
}

// UpdateReminderPosition ...
func UpdateReminderPosition(v Reminder) (*Reminder, error) {
	var getDB = db.GetDB()
	if err := getDB.Where("text = ?", v.Text).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateReminderNotificationTime ...
func UpdateReminderNotificationTime(id int64) (*Reminder, error) {
	var getDB = db.GetDB()
	r := Reminder{}

	now := time.Now()
	if err := getDB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	r.Time = now.Add(10 * time.Minute).Format(time.RFC3339)
	r.Position = 0
	return &r, nil
}
