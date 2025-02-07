package chanselect

import (
	"fmt"

	"github.com/lizthejester/lizbotgo/pkg/alarm"
)

type Server struct {
	MainChannel         string
	ExpiredAlarmManager alarm.AlarmManager
}

type ServerManager struct {
	Servers map[string]*Server
}

func NewServerManager() *ServerManager {
	newServerManager := &ServerManager{
		Servers: make(map[string]*Server),
		/* mapname = make(map[keytype]valuetype) */
	}
	return newServerManager
}

func (m *ServerManager) GetServer(GuildID string) *Server {
	if m.Servers == nil {
		fmt.Println("map is empty!")
	}
	_, found := m.Servers[GuildID]
	if found {
		return m.Servers[GuildID]
	}
	newServer := &Server{
		MainChannel: "",
	}
	//m.Servers = append(m.Servers, *newServer)
	m.Servers[GuildID] = newServer
	return newServer
}

func (s *Server) SetChannel(ChannelID string) {
	s.MainChannel = ChannelID
}

func (m *ServerManager) InitServer(serverid string, mainchannel string) {
	m.Servers[serverid] = &Server{
		MainChannel: mainchannel,
	}
}
