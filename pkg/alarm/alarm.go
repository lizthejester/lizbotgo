package alarm

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type AlarmManager struct {
	Alarms []Alarm
}
type Alarm struct {
	ChannelID string
	Deadline  time.Time
	Content   string
	Name      string
	UserID    string
}

func (m *AlarmManager) GetAlarms() []Alarm {
	return m.Alarms
}

func (m *AlarmManager) SetAlarm(alarm *Alarm) {
	var err error
	m.Alarms = append(m.Alarms, *alarm)
	alarmIndex := len(m.Alarms)
	timer1 := time.NewTimer(time.Until(alarm.Deadline))
	<-timer1.C
	fmt.Println("timer fired")
	m.Alarms[alarmIndex].Deadline, err = time.Parse("01 02 2006 03:04PM -0700", "01 02 2006 03:04PM -0700")
	if err != nil {
		fmt.Println(err)
	}
}
