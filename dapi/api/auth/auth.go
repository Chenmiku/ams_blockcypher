package auth

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/user"
	"net/http"
	"strings"
)

type AuthServer struct {
	*http.ServeMux
	web.JsonServer
}

func NewAuthServer() *AuthServer {
	var s = &AuthServer{
		ServeMux: http.NewServeMux(),
	}
	s.HandleFunc("/login", s.HandleLogin)
	s.HandleFunc("/get_profile", s.HandleGetProfile)
	s.HandleFunc("/logout", s.HandleLogout)
	s.HandleFunc("/change_pass", s.handleChangePass)
	return s
}

func (s *AuthServer) MustGetUser(r *http.Request) (*user.User, error) {
	var id = session.MustGet(r).UserID
	var v, err = user.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *AuthServer) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	u, err := s.MustGetUser(r)
	if user.TableUser.IsErrNotFound(err) {
		s.ErrorMessage(w, "user_not_found")
		return
	}

	s.SendDataSuccess(w, map[string]interface{}{
		"user": &u,
	})
}

func (s *AuthServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body = struct {
		Email    string
		Password string
	}{}

	s.MustDecodeBody(r, &body)

	var u, err = user.GetByEmail(strings.ToLower(body.Email))
	if user.TableUser.IsErrNotFound(err) {
		s.SendError(w, err) //(w, "user_not_found")
		return
	}
	web.AssertNil(err)

	if err = u.ComparePassword(body.Password); err != nil {
		s.SendError(w, err) //(w, "password_not_campare")
		return
	}

	var ses = session.MustNew(u)
	s.SendData(w, map[string]interface{}{
		"user":  u,
		"token": ses.ID,
	})
}

func (s *AuthServer) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session.MustClear(r)
	s.SendData(w, nil)
}

func (s *AuthServer) handleChangePass(w http.ResponseWriter, r *http.Request) {
	var body = struct {
		OldPass   string `json:"old_pass"`
		NewPass   string `json:"new_pass"`
		ReNewPass string `json:"re_new_pass"`
		Email     string `json:"email"`
	}{}

	s.MustDecodeBody(r, &body)

	var u, err = user.GetByEmail(strings.ToLower(body.Email))
	if user.TableUser.IsErrNotFound(err) {
		s.ErrorMessage(w, "user_not_found")
		return
	}

	if err := u.ComparePassword(body.OldPass); err != nil {
		s.ErrorMessage(w, "password_not_campare")
		return
	}
	err = u.UpdatePass(body.NewPass)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendData(w, map[string]interface{}{
			"status": "success",
		})
	}
}
