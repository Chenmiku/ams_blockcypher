package test

import (
	"http/web"
	"net/http"
)

type TestServer struct {
	*http.ServeMux
	web.JsonServer
}

func NewTestServer() *TestServer {

	var s = &TestServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/join", s.Join)
	return s
}

func (t *TestServer) Join(w http.ResponseWriter, r *http.Request) {
	// s := socket.NewClient(w, r, nil)
	// s.Recived = func(data []byte) {
	// 	s.Send(data)
	// }

	// s.Closed = func(code int, text string) {
	// 	fmt.Println(fmt.Sprintf("%v %v", code, text))
	// }

}
