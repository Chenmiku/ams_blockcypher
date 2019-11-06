package activity

import (
	"db/mgo"
	"http/web"
	"myproject/dapi/o/model"
	"time"
)

const errMissingData = web.BadRequest("initialize activity without data")

var ActivityTable = model.NewUnsafeTable("activities", "act")

func (a *Activity) Create(data interface{}) error {
	if data == nil {
		return errMissingData
	}
	now := time.Now()
	a.StartAt = now.Unix()
	a.Date = now.Format(dateFormat)
	a.EndAt = a.StartAt
	a.ActiveDuration = 0
	a.Data = data
	return ActivityTable.UnsafeCreate(a)
}

func (a *Activity) Ping(data interface{}) error {
	a.EndAt = time.Now().Unix()
	a.ActiveDuration = a.EndAt - a.StartAt
	update := mgo.M{
		"e_at": a.EndAt,
		"a_d":  a.ActiveDuration,
	}
	if data != nil {
		update["data"] = data
		a.Data = data
	}
	return ActivityTable.UnsafeUpdateByID(a.ID, update)
}
