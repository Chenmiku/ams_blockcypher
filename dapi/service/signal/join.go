package signal

import (
	"fmt"
	"log"
	"ams_system/dapi/o/report/performance"
	"net/http"
	"socket"
	"time"
)

func (s *SignalServer) mustGet(w http.ResponseWriter, r *http.Request) *aud.Device {
	var dtoken = r.URL.Query().Get("dtoken")
	dtok, err := aud.GetActiveDeviceByID(dtoken)

	if err != nil {
		s.SendJson(w, map[string]interface{}{
			"status": "error",
			"error":  "dtoken not found or deactivate",
		})
	}

	return dtok
}

func (s *SignalServer) Join2(w http.ResponseWriter, r *http.Request) {

	dtoken := s.mustGet(w, r)

	//report
	var performance = &performance.Performance{}
	performance.DeviceID = dtoken.DeviceID
	performance.UserID = s.Session(w, r).UserID
	performance.OnTime = 0

	onJoin := func(client *socket.Client) error {

		// client.On("run_campain", func(args ...interface{}) {

		// 	var data = fmt.Sprintf("%v", args[0])

		// 	c, err := campain.GetByID(strings.Trim(data, " "))

		// 	performance.Period.Start = c.TimeLive.Start
		// 	performance.Period.End = c.TimeLive.End

		// 	if err != nil {
		// 		client.SendError("run_campain", err)
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

		// 		client.SendSuccess("run_campain")

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
		// })

		//check device
		d, err := device.GetByID(dtoken.DeviceID)

		if err != nil {
			return err
		}

		performance.BranchID = d.BranchID

		//check duplicate device id
		devices := Get(dtoken.DeviceID)

		if len(devices) > 0 {
			return fmt.Errorf("device %v already connected", dtoken.DeviceID)
		}

		Add(client, dtoken.DeviceID, s.Session(w, r))

		return nil
	}

	onLeave := func(client *socket.Client) {

		d := FindDeviceBySocketID(client.ID)

		Remove(d)

		performance.LogoutTime = time.Now().Unix()
		performance.OnTime = performance.LogoutTime - performance.CTime.Unix()

		err := performance.Update(performance)
		if err != nil {
			log.Printf("performance update error: %v", err)
		}
	}

	go socket.NewClient(socket.Upgrade(w, r, nil), onJoin, onLeave)

	err := performance.Create()

	if err != nil {
		log.Printf("performance create error: %v", err)
	}
}
