package library

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"net/http"
	"path"
)

type FileUploadServer struct {
	web.JsonServer
	*http.ServeMux
	Dir   string
	Limit int64
}

func (f *FileUploadServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {
	ses, err := session.FromContext(r.Context())

	if err != nil {
		f.SendError(w, err)
	}

	return ses
}

func NewFileUploadServer(dir string, limit int64) *FileUploadServer {
	var s = &FileUploadServer{
		ServeMux: http.NewServeMux(),
		Dir:      dir,
		Limit:    limit,
	}

	s.HandleFunc("/get_by_chapter", s.HandleGetByChapter)
	s.HandleFunc("/upload", s.HandleUploadFile)
	s.HandleFunc("/get_by_name", s.HandleGetByName)
	s.HandleFunc("/update", s.HandleUpdateFile)
	s.HandleFunc("/mark_delete", s.HandleMarkDeleteFile)
	return s
}

func (s *FileUploadServer) requestFolder(r *http.Request) string {
	var folder = r.FormValue("folder")
	if len(folder) < 2 {
		folder = path.Clean(r.URL.Path)
	}
	return folder
}

func (s *FileUploadServer) filename(r *http.Request) string {
	return r.FormValue("name")
}

func (s *FileUploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.Recover(w)
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add(
		"Access-Control-Allow-Methods",
		"OPTIONS, HEAD, GET, POST, PUT, DELETE",
	)
	header.Add(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
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
		panic(web.InternalServerError("session_context_error"))
	}

	newRequest := r.WithContext(newContext)

	s.ServeMux.ServeHTTP(w, newRequest)
}
