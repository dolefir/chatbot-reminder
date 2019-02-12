package services

import (
	"github.com/mileusna/facebook-messenger"
	"github.com/simplewayua/chatbot-reminder/models"
	"strconv"
	"strings"
	"time"
)

// MonitorBotNotification notification
func MonitorBotNotification(msng *messenger.Messenger) {
	for {
		reminders, err := models.GetAllTimeReminder()
		if err != nil {
			continue
		}
		t := time.Now()
		fTimeNow := formatTimeNow(&t)

		for _, v := range reminders {
			fTimeDB := dbFormatTime(v.Time)

			if fTimeDB == fTimeNow && v.Position == 0 {
				nameToInt, _ := strconv.ParseInt(v.NameID, 10, 64)
				gm := PostBackMessage(msng, nameToInt, v)
				msng.SendMessage(gm)

				reminder, err := models.UpdateReminderPosition(v)
				if err != nil {
					continue
				}
				reminder.Position = 1
				reminder.SaveData()
			}
		}
		time.Sleep(time.Minute)
	}
}

// PostBackMessage ...
func PostBackMessage(msng *messenger.Messenger, u int64, r models.Reminder) (g *messenger.GenericMessage) {
	btn1 := msng.NewPostbackButton("accept", "ACCEPT")
	btn2 := msng.NewPostbackButton("snooze", "SNOOZE")
	gm := msng.NewGenericMessage(u)
	gm.AddElement(messenger.Element{Title: r.Text, Subtitle: "Notification", Buttons: []messenger.Button{btn1, btn2}})
	return &gm
}

func formatTimeNow(t *time.Time) (timeF string) {
	tm := strings.SplitAfterN(t.Format(time.RFC3339), ":", 3)
	timeF = tm[0] + tm[1]
	return
}

func dbFormatTime(t string) (timeF string) {
	tm := strings.SplitAfterN(t, ":", 3)
	timeF = tm[0] + tm[1]
	return
}
