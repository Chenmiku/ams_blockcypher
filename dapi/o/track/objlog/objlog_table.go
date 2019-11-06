package objlog

import (
	"myproject/dapi/o/model"
)

var TableObjLog = model.NewUnsafeTable("objlog", "obj")

// Save the obj to storage backend
func (o *ObjLog) Save(auth string, ip string) error {
	o.Author = auth
	o.IP = ip
	return TableObjLog.UnsafeCreate(o)
}
