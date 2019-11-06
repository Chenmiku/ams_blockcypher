package public

import (
	"http/web"
	"myproject/dapi/api/public/category"
	"myproject/dapi/api/public/chapter"
	"myproject/dapi/api/public/comic"
	"myproject/dapi/api/public/menu"
	"myproject/dapi/api/public/org"
	"myproject/dapi/config"
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
	s.Handle("/category/", http.StripPrefix("/category", category.NewCategoryServer()))
	s.Handle("/chapter/", http.StripPrefix("/chapter", chapter.NewChapterServer()))
	s.Handle("/comic/", http.StripPrefix("/comic", comic.NewComicServer()))
	s.Handle("/menu/", http.StripPrefix("/menu", menu.NewMenuServer()))
	return s
}
