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

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
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

func (s *MenuServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var v = &menu.Menu{}
	s.MustDecodeBody(r, v)
	err := v.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, v)
	}
}

func (s *MenuServer) mustGetMenu(r *http.Request) (*menu.Menu, error) {
	var id = r.URL.Query().Get("id")
	var v, err = menu.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *MenuServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newMenu = &menu.Menu{}
	s.MustDecodeBody(r, newMenu)
	v, err := s.mustGetMenu(r)
	if err != nil {
		s.ErrorMessage(w, "menu_not_found")
		return
	}
	err = v.Update(newMenu)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := menu.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "menu_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func (s *MenuServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetMenu(r)
	if err != nil {
		s.ErrorMessage(w, "menu_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}

func (s *MenuServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetMenu(r)
	if err != nil {
		s.ErrorMessage(w, "menu_not_found")
		return
	}
	err = menu.MarkDelete(v.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
