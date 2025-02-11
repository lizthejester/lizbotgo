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
}

func (m *AlarmManager) GetAlarms() []Alarm {
	return m.Alarms
}

func (m *AlarmManager) SetAlarm(alarm *Alarm, session *discordgo.Session, channelid string) {
	fmt.Println("Attempting to set alarm", alarm.Name)
	var err error
	m.Alarms = append(m.Alarms, *alarm)
	alarmIndex := len(m.Alarms) - 1
	deadlineTime, err := time.Parse("01 02 2006 03:04PM -0700", alarm.Deadline)
	if err != nil {
		fmt.Println(err)
	}
	timer1 := time.NewTimer(time.Until(deadlineTime))
	<-timer1.C
	fmt.Println("timer fired")
	if m.Alarms[alarmIndex].Deadline == "01 02 2006 03:04PM -0700" {
		return
	}
	session.ChannelMessageSend(channelid, alarm.Name+" went off!")
	m.Alarms[alarmIndex].Deadline = "01 02 2006 03:04PM -0700"
}
