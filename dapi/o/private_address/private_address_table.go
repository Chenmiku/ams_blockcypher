package private_address

import (
	"ams_system/dapi/o/model"
	"time"
)

var TablePrivateAddress = model.NewTable("private_addresss")

func (w *PrivateAddress) Create() error {
	w.CTime = time.Now().Unix()

	return TablePrivateAddress.Create(w)
}

func MarkDelete(id string) error {
	return TablePrivateAddress.MarkDelete(id)
}

func (w *PrivateAddress) UpdateById(newvalue *PrivateAddress) error {
	return TablePrivateAddress.UnsafeUpdateByID(w.ID, newvalue)
}
