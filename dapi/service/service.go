package service

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/service/monitor"
	"myproject/dapi/service/signal"
	"myproject/dapi/service/test"
	"net/http"
)

type ServiceServer struct {
	*http.ServeMux
	web.JsonServer
}

func NewServiceServer() *ServiceServer {

	var s = &ServiceServer{
		ServeMux: http.NewServeMux(),
	}

	s.Handle("/signal/", http.StripPrefix("/signal", signal.NewSignalServer()))
	s.Handle("/monitor/", http.StripPrefix("/monitor", monitor.NewMonitorServer()))
	s.Handle("/test/", http.StripPrefix("/test", test.NewTestServer()))

	return s
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.Recover(w)
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	header.Add(
		"Access-Control-Allow-Methods",
		"OPTIONS, HEAD, GET, POST, DELETE",
	)
	header.Add(
		"Access-Control-Allow-Headers",
		"Content-Type, Content-Range, Content-Disposition",
	)
	header.Add(
		"Access-Control-Allow-Credentials",
		"true",
	)
	header.Add(
		"Access-Control-Max-Age",
		"2520000", // 30 days
	)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	newContext, err := session.NewContext(r.Context(), session.MustGet(r))

	if err != nil {
		panic(web.InternalServerError("session context error"))
	}

	newRequest := r.WithContext(newContext)

	s.ServeMux.ServeHTTP(w, newRequest)
}
