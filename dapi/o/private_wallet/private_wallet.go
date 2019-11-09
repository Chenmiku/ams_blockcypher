package private_wallet

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectPrivateWalletLog = mlog.NewTagLog("object_private_wallet")

//PrivateWallet
type PrivateWallet struct {
	mgo.BaseModel `bson:",inline"`
	Address     string `bson:"address,omitempty" json:"address"`
	PublicKey      string `bson:"public_key,omitempty" json:"public_key"`
	PrivateKey      string `bson:"private_key,omitempty" json:"private_key"`
	Wif      string `bson:"wif,omitempty" json:"wif"`
	CoinType      string `bson:"coin_type,omitempty" json:"coin_type"`
	ParentID      string `bson:"parent_id,omitempty" json:"parent_id"`
	UserID        string `bson:"user_id,omitempty" json:"user_id"`
	CTime         int64  `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanPrivateWallet() interface{} {
	return &PrivateWallet{}
}
