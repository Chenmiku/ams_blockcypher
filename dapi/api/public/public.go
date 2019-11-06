package public

import (
	"ams_system/dapi/api/public/org"
	"ams_system/dapi/config"
	"http/web"
	"net/http"
)

type PublicServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewPublicServer(pc *config.ProjectConfig) *PublicServer {
	var s = &PublicServer{
		ServeMux: http.NewServeMux(),
	}

	s.Handle("/org/", http.StripPrefix("/org", org.NewOrgServer()))
	return s
}
