package models

import (
	"encoding/json"
	"github.com/simplewayua/chatbot-reminder/db"
	"github.com/simplewayua/chatbot-reminder/dialogflowmap"
)

// Reminder ...
type Reminder struct {
	ID   uint   `gorm:"primary_key"`
	Text string `gorm:"NOT NULL"`
	Time string
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
		Text: npl.Entities.Text,
		Time: npl.Entities.Datewithtime,
	}
	return &reminder, nil
}
