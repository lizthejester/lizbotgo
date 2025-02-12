package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lizthejester/lizbotgo/pkg/alarm"
	"github.com/lizthejester/lizbotgo/pkg/chanselect"
	"github.com/lizthejester/lizbotgo/pkg/deck"
)

func NewManager() *UserManager {
	return &UserManager{
		users: map[string]*User{},
	}
}

type User struct {
	AlarmManager *alarm.AlarmManager
	TarotDeck    *deck.DeckManager
	UserID       string
}

type UserManager struct {
	users map[string]*User
}

func (m *UserManager) GetAllUsers() map[string]*User {
	return m.users
}

func (m *UserManager) GetUser(id string, s *discordgo.Session, sm *chanselect.ServerManager) *User {
	user, exists := m.users[id]
	if exists {
		return user
	}

	newUser := &User{
		TarotDeck:    deck.InitDeck(),
		AlarmManager: m.InitAlarms(id, s, sm),
		UserID:       id,
	}

	m.users[id] = newUser

	return newUser
}

func (m *UserManager) InitAlarms(userid string, session *discordgo.Session, sm *chanselect.ServerManager) *alarm.AlarmManager {
	//var alarms []alarm.Alarm
	alarms := []alarm.Alarm{}
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=10000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from alarms where userid = ?", userid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var newalarm alarm.Alarm
		var dbid int
		var deadlineString string
		if err = rows.Scan(&dbid, &newalarm.Name, &deadlineString, &newalarm.Content, &newalarm.ChannelID, &newalarm.UserID, &newalarm.ServerID, &newalarm.LoopFreq); err != nil {
			fmt.Println(err, dbid)
		}
		newalarm.Deadline = deadlineString
		if err != nil {
			fmt.Println(err)
		}
		alarms = append(alarms, newalarm)
	}
	newAlarmManager := &alarm.AlarmManager{
		Alarms: []alarm.Alarm{},
	}
	for _, v := range alarms {
		deadline, err := time.Parse("01 02 2006 03:04PM -0700", v.Deadline)
		if err != nil {
			fmt.Println(err)
		}
		if time.Until(deadline) < 0 {
			sm.GetServer(v.ServerID).ExpiredAlarmManager.Alarms = append(sm.Servers[v.ServerID].ExpiredAlarmManager.Alarms, v)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			go newAlarmManager.SetAlarm(&v, session, v.ChannelID)
		}
	}
	return newAlarmManager
}
