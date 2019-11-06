package chapter

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/chapter"
	"net/http"
	"strconv"
)

type ChapterServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewChapterServer() *ChapterServer {
	var s = &ChapterServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/get_all", s.mustGetChapterAll)
	s.HandleFunc("/get_by_comic", s.mustGetChapterAllByComicID)
	return s
}

func (s *ChapterServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func (s *ChapterServer) mustGetChapterAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []chapter.Chapter{}

	count, err := chapter.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"chapters": res,
			"count":    count,
		})
	}
}

func (s *ChapterServer) mustGetChapterAllByComicID(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	id := r.URL.Query().Get("comic_id")

	var res = []chapter.Chapter{}

	count, err := chapter.GetByComicID(pageSize, pageNumber, sortBy, sortOrder, &res, id)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"chapters": res,
			"count":    count,
		})
	}
}

func (s *ChapterServer) mustGetChapter(r *http.Request) (*chapter.Chapter, error) {
	var id = r.URL.Query().Get("id")
	var v, err = chapter.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *ChapterServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetChapter(r)
	if err != nil {
		s.ErrorMessage(w, "chapter_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}
