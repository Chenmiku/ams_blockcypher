package wallet

import (
	"errors"
	"ams_system/dapi/o/model"
	"time"
)

var TableWallet = model.NewTable("wallets")

func (w *Wallet) Create() error {
	if err := w.ensureUniqueName(); err != nil {
		return errors.New("name_already_exists")
	}

	w.CTime = time.Now().Unix()

	return TableWallet.Create(w)
}

func MarkDelete(id string) error {
	return TableWallet.MarkDelete(id)
}

func (w *Wallet) UpdateById(newvalue *Wallet) error {
	if newvalue.Name != w.Name {
		if err := newvalue.ensureUniqueName(); err != nil {
			return errors.New("name_already_exists")
		}
	}

	return TableWallet.UnsafeUpdateByID(w.ID, newvalue)
}


