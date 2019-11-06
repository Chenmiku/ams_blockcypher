package signal

import (
	"fmt"
	"ams_system/dapi/api/auth/session"
	"socket"
	"time"
)

type Device struct {
	DeviceID string
	Socket   *socket.Client
	Session  *session.Session
	CTime    time.Time
}

var Devices = []*Device{}

type changed = struct {
	UserID  string
	Devices []*Device
}

var Changed = make(chan changed, 102400)

type updateRunCampain = struct {
	UserID     string
	Campain    string
	DeviceID   string
	OnlineTime time.Time
}

var UpdateRunCampain = make(chan updateRunCampain, 102400)

func Add(socket *socket.Client, deviceID string, session *session.Session) {

	Devices = append(Devices, &Device{
		Socket:   socket,
		DeviceID: deviceID,
		Session:  session,
		CTime:    time.Now(),
	})

	c := changed{
		UserID:  session.UserID,
		Devices: ListIDByUser(session),
	}

	select {
	case Changed <- c:
	default:
	}
}

func FindDeviceBySocketID(id string) *Device {
	for _, d := range Devices {
		if d.Socket.ID == id {
			return d
		}
	}
	return nil
}

func Remove(device *Device) {
	for i, d := range Devices {
		if d.Socket.ID == device.Socket.ID {

			Devices = append(Devices[:i], Devices[i+1:]...)

			c := changed{
				UserID:  d.Session.UserID,
				Devices: ListIDByUser(d.Session),
			}

			select {
			case Changed <- c:
			default:
			}
		}
	}

}

func Get(deviceID string) []*Device {

	var temp []*Device

	for _, d := range Devices {
		if d.DeviceID == deviceID {
			temp = append(temp, d)
		}
	}

	return temp
}

func List() {
	for _, d := range Devices {
		fmt.Println("Signal User ID:", d.Session.UserID)
	}
}

func ListIDByUser(session *session.Session) []*Device {

	temp := []*Device{}

	for _, d := range Devices {
		if d.Session.UserID == session.UserID {
			temp = append(temp, d)
		}
	}

	return temp
}
