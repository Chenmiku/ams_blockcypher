package monitor

import (
	"http/web"
	"log"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/report/performance"
	"myproject/dapi/service/signal"
	"net/http"
	"socket"
)

type MonitorServer struct {
	*http.ServeMux
	web.JsonServer
}

func NewMonitorServer() *MonitorServer {

	var s = &MonitorServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/join", s.Join)
	return s
}

func (s *MonitorServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

func (s *MonitorServer) Join(w http.ResponseWriter, r *http.Request) {

	onJoin := func(client *socket.Client) error {

		client.Recived = func(uri string, data string) {

		}

		client.SendJson("/devices", buildMonitorData(client, s.Session(w, r)))

		Add(client, s.Session(w, r))

		return nil
	}

	onLeave := func(client *socket.Client) {
		m := FindMonitorBySocketID(client.ID)
		Remove(m)
	}

	socket.NewClient(socket.Upgrade(w, r, nil), onJoin, onLeave)
}

func buildMonitorData(socket *socket.Client, session *session.Session) []*MonitorData {

	var mdata = []*MonitorData{}

	devices := signal.ListIDByUser(session)

	for _, d := range devices {

		var m = &MonitorData{}

		var p, err = performance.GetByLatest(d.DeviceID)

		if err == nil {
			m.CampainID = p.CampainID
		} else {
			log.Println(err)
		}

		m.DeviceID = d.DeviceID
		m.OnlineTime = d.CTime

		mdata = append(mdata, m)
	}

	return mdata
}

func init() {
	go func() {
		for {
			//tracking signal changed
			select {
			case deviceChanged := <-signal.Changed:

				monitors := Get(deviceChanged.UserID)

				for _, m := range monitors {
					m.Socket.SendJson("/devices", buildMonitorData(m.Socket, m.Session))
				}

			case campainChanged := <-signal.UpdateRunCampain:

				var mdata = MonitorData{}

				mdata.CampainID = campainChanged.Campain
				mdata.DeviceID = campainChanged.DeviceID
				mdata.OnlineTime = campainChanged.OnlineTime

				monitors := Get(campainChanged.UserID)

				for _, m := range monitors {
					m.Socket.SendJson("/campain_update", mdata)
				}
			}

		}
	}()
}
