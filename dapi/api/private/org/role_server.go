package org

import (
	"gopkg.in/mgo.v2/bson"
	"http/web"
	"ams_system/dapi/o/org/role"
	"net/http"
)

type RoleServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewRoleServer() *RoleServer {
	var s = &RoleServer{
		ServeMux: http.NewServeMux(),
	}
	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	s.HandleFunc("/get_all", s.HandleAllRole)
	return s
}

func (s *RoleServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var u = &role.Role{}
	s.MustDecodeBody(r, u)
	web.AssertNil(u.Create())
	s.SendData(w, u)
}

func (s *RoleServer) mustGetRole(r *http.Request) *role.Role {
	var id = r.URL.Query().Get("id")
	var u, err = role.GetByID(id)
	web.AssertNil(err)
	return u
}

func (s *RoleServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newRole = &role.Role{}
	s.MustDecodeBody(r, newRole)
	var u = s.mustGetRole(r)
	web.AssertNil(u.Update(newRole))
	s.SendData(w, nil)
}

func (s *RoleServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	var u = s.mustGetRole(r)
	s.SendData(w, u)
}

func (s *RoleServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	var u = s.mustGetRole(r)
	web.AssertNil(role.MarkDelete(u.ID))
	s.Success(w)
}

func (s *RoleServer) HandleAllRole(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []role.Role{}

	where := bson.M{
		"dtime": 0,
	}

	count, err := role.RoleTable.ReadPagingSort(where, pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendData(w, map[string]interface{}{
			"roles": res,
			"count": count,
		})
	}
}
