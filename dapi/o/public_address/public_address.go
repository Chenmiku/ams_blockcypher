package public_address

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectPrivateAddressLog = mlog.NewTagLog("object_public_address")

//PublicAddress
type PublicAddress struct {
	mgo.BaseModel          `bson:",inline"`
	Address                string `bson:"address,omitempty" json:"address"`
	WalletID               string `bson:"wallet_id,omitempty" json:"wallet_id"`
	TotalRevceived         int    `bson:"total_revceived,omitempty" json:"total_revceived"`
	TotalSent              int    `bson:"total_sent,omitempty" json:"total_sent"`
	Balance                int    `bson:"balance,omitempty" json:"balance"`
	UnconfirmedBalance     int    `bson:"unconfirmed_balance,omitempty" json:"unconfirmed_balance"`
	FinalBalance           int    `bson:"final_balance,omitempty" json:"final_balance"`
	CoinType               string `bson:"coin_type,omitempty" json:"coin_type"`
	ConfirmedTransaction   int    `bson:"confirmed_transaction,omitempty" json:"confirmed_transaction"`
	UnconfirmedTransaction int    `bson:"unconfirmed_transaction,omitempty" json:"unconfirmed_transaction"`
	FinalTransaction       int    `bson:"final_transaction,omitempty" json:"final_transaction"`
	CTime                  int64  `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanPublicAddress() interface{} {
	return &PublicAddress{}
}
