package contact

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/contact"
	"net/http"
	"strconv"
)

type ContactServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewContactServer() *ContactServer {
	var s = &ContactServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	s.HandleFunc("/get_all", s.mustGetBranchAll)
	return s
}

func (s *ContactServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

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

func (s *ContactServer) mustGetBranchAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []contact.Contact{}

	count, err := contact.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"contacts": res,
			"count":    count,
		})
	}
}

func (s *ContactServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var v = &contact.Contact{}
	s.MustDecodeBody(r, v)
	err := v.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, v)
	}
}

func (s *ContactServer) mustGetContact(r *http.Request) (*contact.Contact, error) {
	var id = r.URL.Query().Get("id")
	var v, err = contact.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *ContactServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newContact = &contact.Contact{}
	s.MustDecodeBody(r, newContact)
	v, err := s.mustGetContact(r)
	if err != nil {
		s.ErrorMessage(w, "contact_not_found")
		return
	}
	err = v.Update(newContact)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := contact.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "contact_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func (s *ContactServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetContact(r)
	if err != nil {
		s.ErrorMessage(w, "contact_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}

func (s *ContactServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetContact(r)
	if err != nil {
		s.ErrorMessage(w, "contact_not_found")
		return
	}
	err = contact.MarkDelete(v.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
