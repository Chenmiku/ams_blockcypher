package signal

import (
	"fmt"
	"http/web"
	"ams_system/dapi/api/auth/session"
	"net/http"
	"socket"
)

type SignalServer struct {
	*http.ServeMux
	web.JsonServer
}

func NewSignalServer() *SignalServer {

	var s = &SignalServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/join", s.Join)
	s.HandleFunc("/send", s.Send)
	s.HandleFunc("/all", s.AllDeviceHandler)

	return s
}

func (s *SignalServer) AllDeviceHandler(w http.ResponseWriter, r *http.Request) {
	device := []string{}
	for _, d := range Devices {
		device = append(device, d.DeviceID)
	}
	s.SendData(w, device)
}

func (s *SignalServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

func (s *SignalServer) Join(w http.ResponseWriter, r *http.Request) {

	var deviceID = r.URL.Query().Get("device_id")

	// //report
	// var performance = &performance.Performance{}
	// performance.DeviceID = deviceID
	// performance.UserID = s.Session(w, r).UserID
	// performance.OnTime = 0

	onJoin := func(client *socket.Client) error {

		client.Recived = func(uri string, data string) {

			// if uri == "/run_campain" {

			// 	c, err := campain.GetByID(strings.Trim(data, " "))

			// 	performance.Period.Start = c.TimeLive.Start
			// 	performance.Period.End = c.TimeLive.End

			// 	if err != nil {
			// 		client.SendError(uri, err)
			// 	} else {
			// 		if performance.CampainID == "" {
			// 			performance.CampainID = c.ID
			// 			err := performance.Update(performance)
			// 			if err != nil {
			// 				log.Println(err)
			// 			}
			// 		} else if performance.CampainID != c.ID {
			// 			performance.CampainID = c.ID
			// 			err = performance.Create()
			// 			if err != nil {
			// 				log.Println(err)
			// 			}
			// 		}

			// 		client.SendSuccess(uri)

			// 		d := FindDeviceBySocketID(client.ID)

			// 		u := updateRunCampain{
			// 			UserID:     s.Session(w, r).UserID,
			// 			Campain:    c.ID,
			// 			DeviceID:   d.DeviceID,
			// 			OnlineTime: d.CTime,
			// 		}
			// 		select {
			// 		case UpdateRunCampain <- u:
			// 		default:
			// 		}

			// 	}
			// }

		}

		//check device
		_, err := device.GetByID(deviceID)

		if err != nil {
			return err
		}

		// performance.BranchID = d.BranchID

		//check duplicate device id
		devices := Get(deviceID)

		if len(devices) > 0 {
			return fmt.Errorf("device %v already connected", deviceID)
		}

		Add(client, deviceID, s.Session(w, r))

		return nil
	}

	onLeave := func(client *socket.Client) {

		d := FindDeviceBySocketID(client.ID)

		Remove(d)

		// performance.LogoutTime = time.Now().Unix()
		// performance.OnTime = performance.LogoutTime - performance.CTime.Unix()

		// err := performance.Update(performance)
		// if err != nil {
		// 	log.Printf("performance update error: %v", err)
		// }
	}

	go socket.NewClient(socket.Upgrade(w, r, nil), onJoin, onLeave)

	// err := performance.Create()

	// if err != nil {
	// 	log.Printf("performance create error: %v", err)
	// }
}

func (s *SignalServer) Send(w http.ResponseWriter, r *http.Request) {

	var deviceID = r.URL.Query().Get("device_id")

	devices := Get(deviceID)

	if len(devices) != 0 {

		for _, d := range devices {
			d.Socket.SendUri("/reload")
		}

		s.SendData(w, map[string]interface{}{
			"action":    "/reload",
			"device_id": deviceID,
		})

	} else {
		s.SendError(w, web.NotFound(fmt.Sprintf("device %v not connect", deviceID)))
	}
}
