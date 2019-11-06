package objlog

import (
	"db/mgo"
	"time"
)

type ObjLog struct {
	mgo.UnsafeBaseModel `bson:",inline"`
	EntityName          string      `bson:"name" json:"name"` // Name Of Device
	BranchID            string      `bson:"bid" json:"bid"`   //
	EntityID            string      `bson:"eid" json:"eid"`   //
	Category            string      `bson:"cat" json:"cat"`   //
	Action              string      `bson:"act" json:"act"`
	Author              string      `bson:"aut" json:"aut"`
	IP                  string      `bson:"ip" json:"ip"`
	CreatedAt           int64       `bson:"s_at" json:"s_at"` //
	Data                interface{} `bson:"data" json:"data"`
}

// NewObjLog Create a new activity and insert into database
func NewObjLog(
	entityName string, entityID string,
	branchID string, category string,
	action string, data interface{},
) *ObjLog {
	var a = &ObjLog{
		EntityName: entityName,
		EntityID:   entityID,
		BranchID:   branchID,
		Action:     action,
		Category:   category,
		Data:       data,
	}
	a.CreatedAt = time.Now().Unix()
	return a
}
