package private_wallet

import (
	"ams_system/dapi/o/model"
	"time"
)

var TablePrivateWallet = model.NewTable("private_wallets")

func (w *PrivateWallet) Create() error {
	w.CTime = time.Now().Unix()

	return TablePrivateWallet.Create(w)
}

func MarkDelete(id string) error {
	return TablePrivateWallet.MarkDelete(id)
}

func (w *PrivateWallet) UpdateById(newvalue *PrivateWallet) error {
	return TablePrivateWallet.UnsafeUpdateByID(w.ID, newvalue)
}
