package category

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/category"
	"net/http"
	"strconv"
)

type CategoryServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewCategoryServer() *CategoryServer {
	var s = &CategoryServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.mustGetCategoryAll)
	return s
}

func (s *CategoryServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

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

func (s *CategoryServer) mustGetCategoryAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []category.Category{}

	count, err := category.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"categories": res,
			"count":      count,
		})
	}
}
