package user

import (
	"github.com/lizthejester/lizbotgo/pkg/deck"
)

func NewManager() *UserManager {
	return &UserManager{
		users: map[string]*User{},
	}
}

type User struct {
	Alarms    []Alarm
	TarotDeck *deck.DeckManager
}

type UserManager struct {
	users map[string]*User
}

func (m *UserManager) GetUser(id string) *User {
	user, exists := m.users[id]
	if exists {
		return user
	}

	newUser := &User{
		TarotDeck: deck.InitDeck(),
	}

	m.users[id] = newUser

	return newUser
}
