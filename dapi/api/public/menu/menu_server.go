package menu

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/menu"
	"net/http"
	"strconv"
)

type MenuServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewMenuServer() *MenuServer {
	var s = &MenuServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.HandleGetAll)
	return s
}

func (s *MenuServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

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

func (s *MenuServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []menu.Menu{}

	count, err := menu.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"menus": res,
			"count": count,
		})
	}
}
