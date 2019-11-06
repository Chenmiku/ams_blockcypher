package role

import (
	"ams_system/dapi/o/model"
)

type Role struct {
	model.BaseModel `bson:",inline"`
	Name            string   `bson:"name" json:"name"` //
	Permission      []string `bson:"permission" json:"permission"`
}
