package activity

import (
	"db/mgo"
	"time"
)

type Activity struct {
	mgo.UnsafeBaseModel `bson:",inline"`
	EntityName          string      `bson:"name" json:"name"` // Name Of Device
	BranchID            string      `bson:"bid" json:"bid"`   //
	EntityID            string      `bson:"eid" json:"eid"`   //
	Category            string      `bson:"cat" json:"cat"`   //
	StartAt             int64       `bson:"s_at" json:"s_at"` //
	EndAt               int64       `bson:"e_at" json:"e_at"` //
	ActiveDuration      int64       `bson:"a_d" json:"a_d"`
	Date                string      `bson:"date" json:"date"` //
	Data                interface{} `bson:"data" json:"data"`
}

const dateFormat = "2006-01-02"

// NewActivity Create a new activity and insert into database
func NewActivity(
	entityID string, entityName string,
	branchID string, category string,
	data interface{},
) Activity {
	return Activity{
		EntityName: entityName,
		EntityID:   entityID,
		BranchID:   branchID,
		Category:   category,
		Data:       data,
	}
}

func getDate() string {
	return time.Now().Format(dateFormat)
}
