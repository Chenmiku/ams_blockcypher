package wallet

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectWalletLog = mlog.NewTagLog("object_wallet")

//Wallet
type Wallet struct {
	mgo.BaseModel `bson:",inline"`
	Name          string   `bson:"name,omitempty" json:"name"`
	Token         string   `bson:"token,omitempty" json:"token"`
	Addresses     []string `bson:"addresses,omitempty" json:"addresses,omitempty"`
	UserID        string   `bson:"user_id,omitempty" json:"user_id"`
	CTime         int64    `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanWallet() interface{} {
	return &Wallet{}
}

// check unique name
func (w *Wallet) ensureUniqueName() error {
	if err := TableWallet.NotExist(map[string]interface{}{
		"name": w.Name,
	}); err != nil {
		return err
	}
	return nil
}
