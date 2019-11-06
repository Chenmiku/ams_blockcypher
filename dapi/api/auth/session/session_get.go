package session

import (
	"http/web"
	"myproject/dapi/o/auth/session"
)

const (
	errReadSessonFailed   = web.InternalServerError("read_session_failed")
	errSessionNotFound    = web.Unauthorized("session_not_found")
	errUnauthorizedAccess = web.Unauthorized("unauthorized_access")
)

func Get(sessionID string) (*session.Session, error) {
	var s, err = session.GetByID(sessionID)
	if err != nil {
		if session.TableSession.IsErrNotFound(err) {
			return nil, errSessionNotFound
		}
		sessionLog.Error(err)
		return nil, errReadSessonFailed
	}

	return s, nil
}
