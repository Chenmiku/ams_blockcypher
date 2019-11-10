package public_address

import (
	"ams_system/dapi/o/model"
	"time"
)

var TablePublicAddress = model.NewTable("public_addresses")

func (w *PublicAddress) Create() error {
	w.CTime = time.Now().Unix()

	return TablePublicAddress.Create(w)
}

func MarkDelete(id string) error {
	return TablePublicAddress.MarkDelete(id)
}

func (w *PublicAddress) UpdateById(newvalue *PublicAddress) error {
	return TablePublicAddress.UnsafeUpdateByID(w.ID, newvalue)
}
