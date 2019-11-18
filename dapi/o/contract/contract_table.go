package contract

import (
	"ams_system/dapi/o/model"
	"time"
)

var TableContract = model.NewTable("contracts")

func (w *Contract) Create() error {
	w.CTime = time.Now().Unix()

	return TableContract.Create(w)
}

func MarkDelete(id string) error {
	return TableContract.MarkDelete(id)
}

func (w *Contract) UpdateById(newvalue *Contract) error {
	return TableContract.UnsafeUpdateByID(w.ID, newvalue)
}
