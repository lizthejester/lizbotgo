package alarm

import (
	"fmt"
	"time"
)

type AlarmManager struct {
	Alarms []Alarm
}
type Alarm struct {
	ChannelID string
	Deadline  time.Time
	Content   string
	Name      string
}

func GetAlarms() *AlarmManager {
	return nil
}

func (m *AlarmManager) SetAlarm(alarm *Alarm) {
	timer1 := time.NewTimer(time.Until(alarm.Deadline))
	<-timer1.C
	fmt.Println("timer fired")
}
