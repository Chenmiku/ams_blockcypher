package httpserver

import (
	"ams_system/dapi/api"
	"http/gziphandler"
	"http/static/vstatic"
	"net/http"
	"regexp"
)

func webAssetGzipHandler(handler http.Handler) http.Handler {
	gzip := gziphandler.GzipHandler(handler)
	assetRegex, _ := regexp.Compile(".(js|css)$")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if assetRegex.MatchString(r.URL.Path) {
			gzip.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (phs *ProjectHttpServer) addStaticHandler(s *http.ServeMux) {
	p := phs.pc
	staticConfig := p.Station.Static
	// s.Handle("/", http.RedirectHandler("/app/", http.StatusFound))
	var admin = vstatic.NewVersionStatic(staticConfig.AdminFolder)

	admin.SetUpdate(staticConfig.AppUpdate)
	s.Handle("/", http.StripPrefix("/", webAssetGzipHandler(admin)))

	var device = vstatic.NewVersionStatic(staticConfig.PlayerFolder)
	s.Handle("/player/", http.StripPrefix("/player", webAssetGzipHandler(device)))
}

func (phs *ProjectHttpServer) makeHandler() http.Handler {
	var server = http.NewServeMux()
	phs.addStaticHandler(server)
	// application specific
	apiServer := api.NewApiServer(phs.pc)

	server.Handle("/api/",
		gziphandler.GzipHandler(http.StripPrefix("/api", apiServer)),
	)

	phs.ready <- struct{}{}
	return server
}
