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

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
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

func (s *CategoryServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var v = &category.Category{}
	s.MustDecodeBody(r, v)
	err := v.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, v)
	}
}

func (s *CategoryServer) mustGetCategory(r *http.Request) (*category.Category, error) {
	var id = r.URL.Query().Get("id")
	var v, err = category.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *CategoryServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newCategory = &category.Category{}
	s.MustDecodeBody(r, newCategory)
	v, err := s.mustGetCategory(r)
	if err != nil {
		s.ErrorMessage(w, "category_not_found")
		return
	}
	err = v.Update(newCategory)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := category.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "category_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func (s *CategoryServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetCategory(r)
	if err != nil {
		s.ErrorMessage(w, "category_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}

func (s *CategoryServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetCategory(r)
	if err != nil {
		s.ErrorMessage(w, "category_not_found")
		return
	}
	err = category.MarkDelete(v.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
