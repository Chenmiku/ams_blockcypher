package session

import (
	"ams_system/dapi/o/auth/session"
	"ams_system/dapi/o/org/user"
	"http/web"
	"time"
)

func New(u *user.User) (*session.Session, error) {

	var s = &session.Session{
		UserID: u.ID,
		Email:  u.Email,
		CTime:  time.Now().Unix(),
	}

	var err = s.Create()
	if err != nil {
		sessionLog.Error(err)
		return nil, web.InternalServerError("save_session_failed")
	}
	return s, nil
}

func MustNew(u *user.User) *session.Session {
	s, e := New(u)
	if e != nil {
		panic(e)
	}
	return s
}
