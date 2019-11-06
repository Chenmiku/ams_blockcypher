package private

import (
	"http/web"
	"myproject/dapi/api/private/org"
	"net/http"
)

type orgServer struct {
	web.JsonServer
	*http.ServeMux
}

func newOrgServer() *orgServer {
	var s = &orgServer{
		ServeMux: http.NewServeMux(),
	}
	s.Handle("/user/", http.StripPrefix("/user", org.NewUserServer()))
	s.Handle("/role/", http.StripPrefix("/role", org.NewRoleServer()))
	return s
}
