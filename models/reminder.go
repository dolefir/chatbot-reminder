package models

import (
	"encoding/json"
	"github.com/simplewayua/chatbot-reminder/db"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
	"strings"
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
func (r *Reminder) ToReminder(response *dialogflowmap.NPLResponse) (*Reminder, error) {
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
func (r *Reminder) GetTimesReminder(t string) (rems []Reminder) {
	var getDB = db.GetDB()
	s := strings.Split(t, "T")
	tm := s[0]
	getDB.Where("time LIKE ?", tm+"%").Find(&rems)
	return
}

// DeleteReminder ...
func (r *Reminder) DeleteReminder(t string) {
	var getDB = db.GetDB()
	getDB.Where("text LIKE ?", "%"+t+"%").Delete(Reminder{})
}

// GetAllTimeReminder for notification
func (r *Reminder) GetAllTimeReminder() (rems []Reminder) {
	var getDB = db.GetDB()
	getDB.Find(&rems)
	return rems
}

// UpdateReminderPosition ...
func (r *Reminder) UpdateReminderPosition(v Reminder) *Reminder {
	var getDB = db.GetDB()
	getDB.Where("text = ?", v.Text).First(&v)
	return &v
}
