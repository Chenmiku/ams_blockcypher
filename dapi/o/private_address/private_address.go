package private_address

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectPrivateAddressLog = mlog.NewTagLog("object_private_address")

//PrivateAddress
type PrivateAddress struct {
	mgo.BaseModel `bson:",inline"`
	Address       string `bson:"address,omitempty" json:"address"`
	PublicKey     string `bson:"public_key,omitempty" json:"public_key"`
	PrivateKey    string `bson:"private_key,omitempty" json:"private_key"`
	Wif           string `bson:"wif,omitempty" json:"wif"`
	CoinType      string `bson:"coin_type,omitempty" json:"coin_type"`
	UserID        string `bson:"user_id,omitempty" json:"user_id"`
	ParentID      string `bson:"parent_id,omitempty" json:"parent_id"`
	CTime         int64  `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanPrivateAddress() interface{} {
	return &PrivateAddress{}
}
