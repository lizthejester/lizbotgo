package alarm

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type AlarmManager struct {
	Alarms []Alarm
}
type Alarm struct {
	ChannelID string
	Deadline  string
	Content   string
	Name      string
	UserID    string
	ServerID  string
	LoopFreq  string
}

func (m *AlarmManager) GetAlarms() []Alarm {
	return m.Alarms
}

var SetAlarm func(m *AlarmManager) (alarm *Alarm, session *discordgo.Session, channelid string)

func (m *AlarmManager) SetAlarm(alarm *Alarm, session *discordgo.Session, channelid string) {
	fmt.Println("Attempting to set alarm", alarm.Name)
	var err error
	m.Alarms = append(m.Alarms, *alarm)
	alarmIndex := len(m.Alarms) - 1
	deadlineTime, err := time.Parse("01 02 2006 03:04PM -0700", alarm.Deadline)
	if err != nil {
		fmt.Println(err, "No this one")
	}
	timer1 := time.NewTimer(time.Until(deadlineTime))
	<-timer1.C
	fmt.Println("timer fired")
	if m.Alarms[alarmIndex].Deadline == "01 02 2006 03:04PM -0700" {
		return
	}
	session.ChannelMessageSend(channelid, "<@"+alarm.UserID+"> "+alarm.Name+" went off!\n"+alarm.Content)
	m.Alarms[alarmIndex].Deadline = "01 02 2006 03:04PM -0700"
	switch alarm.LoopFreq {
	case "daily":
		alarm.Deadline = deadlineTime.AddDate(0, 0, 1).Format("01 02 2006 03:04PM -0700")
		m.SetAlarm(alarm, session, channelid)
	case "weekly":
		alarm.Deadline = deadlineTime.AddDate(0, 0, 7).Format("01 02 2006 03:04PM -0700")
		m.SetAlarm(alarm, session, channelid)
	case "monthly":
		alarm.Deadline = deadlineTime.AddDate(0, 1, 0).Format("01 02 2006 03:04PM -0700")
		m.SetAlarm(alarm, session, channelid)
	case "yearly":
		alarm.Deadline = deadlineTime.AddDate(1, 0, 0).Format("01 02 2006 03:04PM -0700")
		m.SetAlarm(alarm, session, channelid)
	default:
		return
	}
}
