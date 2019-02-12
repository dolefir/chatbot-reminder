package services

import (
	"github.com/mileusna/facebook-messenger"
	"github.com/simplewayua/chatbot-reminder/models"
	"strconv"
	"strings"
	"time"
)

var r models.Reminder

// MonitorBotNotification notification
func MonitorBotNotification(msng *messenger.Messenger) {
	for {
		reminders := r.GetAllTimeReminder()
		t := time.Now()
		fTimeNow := formatTimeNow(&t)

		for _, v := range reminders {
			fTimeDB := dbFormatTime(v.Time)

			if fTimeDB == fTimeNow && v.Position == 0 {
				uNameToInt, _ := strconv.ParseInt(v.NameID, 10, 64)
				msng.SendTextMessage(uNameToInt, v.Text)
				reminder := r.UpdateReminderPosition(v)
				reminder.Position = 1
				reminder.SaveData()
			}
		}
		time.Sleep(time.Minute * 2)
	}
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
