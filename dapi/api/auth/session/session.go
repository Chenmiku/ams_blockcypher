package session

import (
	"encoding/json"
	"myproject/dapi/o/org/user"
	"myproject/dapi/x/math"
)

var idMaker = math.RandStringMaker{Length: 40, Prefix: "s"}

type Session struct {
	SessionID string    `json:"id"`
	Email     string    `json:"email"`
	UserID    string    `json:"user_id"`
	Role      user.Role `json:"role"`
	CTime     int64     `json:"ctime"`
}

func (a *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
