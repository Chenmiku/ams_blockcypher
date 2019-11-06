package org

import (
	"ams_system/dapi/o/org/user"
	"http/web"
	"net/http"
	"strconv"
	"strings"
)

type userAPI struct {
	web.JsonServer
	*http.ServeMux
}

func newPublicUserAPI() *userAPI {
	u := new(userAPI)
	u.ServeMux = http.NewServeMux()
	u.HandleFunc("/", u.handleCreate)
	return u
}

func (uapi *userAPI) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	var u = &user.User{}
	uapi.MustDecodeBody(r, u)
	u.Email = strings.ToLower(u.Email)
	err := u.Create()
	if err != nil {
		uapi.ErrorMessage(w, err.Error())
	} else {
		uapi.SendData(w, u)
	}
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
