package contact

import (
	"db/mgo"
)

type Contact struct {
	mgo.BaseModel `bson:",inline"`
	Name          *string  `bson:"name,omitempty" json:"name"`
	Phone         *string  `bson:"phone,omitempty" json:"phone"`
	Email         *string  `bson:"email,omitempty" json:"email"`
	Address       *string  `bson:"address,omitempty" json:"address"`
	Website       *string  `bson:"website,omitempty" json:"website"`
	Lat           *float64 `bson:"lat,omitempty" json:"lat"`
	Lng           *float64 `bson:"lng,omitempty" json:"lng"`
	CTime         int64    `bson:"ctime,omitempty" json:"ctime"` // Create Time
}
