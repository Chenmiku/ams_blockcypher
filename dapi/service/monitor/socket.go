package monitor

import (
	"fmt"
	"myproject/dapi/api/auth/session"
	"socket"
	"time"
)

type Monitor struct {
	Socket  *socket.Client
	Session *session.Session
	CTime   time.Time
}

var Monitors = []*Monitor{}

func Add(socket *socket.Client, session *session.Session) {

	Monitors = append(Monitors, &Monitor{
		Socket:  socket,
		Session: session,
		CTime:   time.Now(),
	})
}

func FindMonitorBySocketID(id string) *Monitor {
	for _, d := range Monitors {
		if d.Socket.ID == id {
			return d
		}
	}
	return nil
}

func FindMonitorByUserID(id string) *Monitor {
	for _, d := range Monitors {
		if d.Session.UserID == id {
			return d
		}
	}
	return nil
}

func Remove(monitor *Monitor) {
	for i, d := range Monitors {
		if d.Socket.ID == monitor.Socket.ID {
			Monitors = append(Monitors[:i], Monitors[i+1:]...)
		}
	}
}

func Get(userID string) []*Monitor {

	var temp = []*Monitor{}

	for _, d := range Monitors {
		if d.Session.UserID == userID {
			temp = append(temp, d)
		}
	}

	return temp
}

func List() {
	for _, h := range Monitors {
		fmt.Println("Monitor UserID:", h.Session.UserID)
	}
}
