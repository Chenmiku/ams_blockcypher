package httpserver

import (
	"http/gziphandler"
	"http/static/vstatic"
	"myproject/dapi/api"
	"myproject/dapi/library"
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

	var seller = vstatic.NewVersionStatic(staticConfig.SellerFolder)
	s.Handle("/seller/", http.StripPrefix("/seller", webAssetGzipHandler(seller)))

	// storage
	storageConfig := p.Station.Storage

	var library = library.NewFileUploadServer(storageConfig.Upload, 0)

	s.Handle("/library/", http.StripPrefix("/library", library))

	s.Handle("/upload/", http.StripPrefix("/upload", http.FileServer(http.Dir(storageConfig.Upload))))
}

func (phs *ProjectHttpServer) makeHandler() http.Handler {
	var server = http.NewServeMux()
	phs.addStaticHandler(server)
	// application specific
	apiServer := api.NewApiServer(phs.pc)
	// service := service.NewServiceServer()

	server.Handle("/api/",
		gziphandler.GzipHandler(http.StripPrefix("/api", apiServer)),
	)

	// server.Handle("/service/",
	// 	http.StripPrefix("/service", service),
	// )

	phs.ready <- struct{}{}
	return server
}
