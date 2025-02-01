package user

import (
	"database/sql"
	"fmt"

	"github.com/lizthejester/lizbotgo/pkg/alarm"
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

func (m *UserManager) GetUser(id string) *User {
	user, exists := m.users[id]
	if exists {
		return user
	}

	newUser := &User{
		TarotDeck:    deck.InitDeck(),
		AlarmManager: m.InitAlarms(id),
		UserID:       id,
	}

	m.users[id] = newUser

	return newUser
}

// name, time, comment, channelid, userid
func (m *UserManager) InitAlarms(userid string) *alarm.AlarmManager {
	var alarms []alarm.Alarm
	db, err := sql.Open("sqlite3", "./lizbot.db")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("select * from alarms where userid = ?", userid)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var newalarm alarm.Alarm
		var dbid int
		if err = rows.Scan(&dbid, &newalarm.Name, &newalarm.Deadline, &newalarm.Content, &newalarm.ChannelID, &newalarm.UserID); err != nil {
			fmt.Println(err, dbid)
		}
		alarms = append(alarms, newalarm)
	}
	newAlarmManager := &alarm.AlarmManager{
		Alarms: alarms,
	}
	return newAlarmManager
}
