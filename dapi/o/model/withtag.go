package model

import (
	"db/mgo"
)

type WithType struct {
	mgo.BaseModel `bson:",inline"`
	Type          string `bson:"type" json:"type"`
}

type IWithType interface {
	mgo.IModel
}

func (t *TableWithType) GetByType(types []string, ptr interface{}) error {
	return t.ReadManyIn("type", types, ptr)
}
