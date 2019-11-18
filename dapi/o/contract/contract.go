package contract

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectContractLog = mlog.NewTagLog("object_contract")

//Contract
type Contract struct {
	mgo.BaseModel   `bson:",inline"`
	Solidity        string   `bson:"solidity,omitempty" json:"solidity"`
	Params          []string `bson:"params,omitempty" json:"params"`
	Publish         []string `bson:"publish,omitempty" json:"publish"`
	PrivateKey      string   `bson:"private_key,omitempty" json:"private_key"`
	GasLimit        int      `bson:"gas_limit,omitempty" json:"gas_limit"`
	Value           int      `bson:"value,omitempty" json:"value"`
	Name            string   `bson:"name,omitempty" json:"name"`
	Bin             string   `bson:"bin,omitempty" json:"bin"`
	Address         string   `bson:"address,omitempty" json:"address"`
	ConfirmedTime   int64    `bson:"confirmed_time,omitempty" json:"confirmed_time"`
	TransactionHash string   `bson:"transaction_hash,omitempty" json:"transaction_hash"`
	Results         []string `bson:"results,omitempty" json:"results"`
	CTime           int64    `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanContract() interface{} {
	return &Contract{}
}
