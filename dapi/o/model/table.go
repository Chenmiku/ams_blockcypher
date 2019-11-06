package model

import (
	"db/mgo"
	"ams_system/dapi/config/cons"
)

type TableWithBranchCode struct {
	*mgo.Table
}

type TableWithCode struct {
	*mgo.Table
}

type TableWithType struct {
	*mgo.Table
}

func NewTable(name string) *mgo.Table {
	var db = GetDB()
	return mgo.NewTable(db, name)
}

func GetDB() *mgo.Database {
	return mgo.GetDB(cons.DB_ID)
}
