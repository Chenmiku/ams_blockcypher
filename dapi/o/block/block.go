package block

import (
	"ams_system/dapi/x/mlog"
	"db/mgo"
)

var objectBlockLog = mlog.NewTagLog("object_block")

//Block
type Block struct {
	mgo.BaseModel          `bson:",inline"`
	Hash                string `bson:"hash,omitempty" json:"hash"`
	Name               string `bson:"name,omitempty" json:"name"`
	Height         int    `bson:"height,omitempty" json:"height"`
	PreviousBlock              string    `bson:"previous_block,omitempty" json:"previous_block"`
	PreviousURL                string    `bson:"previous_url,omitempty" json:"previous_url"`
	Depth     int    `bson:"depth,omitempty" json:"depth"`
	Fees           int    `bson:"fees,omitempty" json:"fees"`
	Version               int `bson:"version,omitempty" json:"version`
	Bits   int    `bson:"bits,omitempty" json:"bits"`
	Nonce int    `bson:"nonce,omitempty" json:"nonce"`
	Transactions       []string    `bson:"transactions,omitempty" json:"transactions"`
	TransactionURL       string    `bson:"transaction_url,omitempty" json:"transaction_url"`
	TotalTransaction       int    `bson:"total_transaction,omitempty" json:"total_transaction"`
	Total       int    `bson:"total,omitempty" json:"total"`
	BlockChainID       string    `bson:"blockchain_id,omitempty" json:"blockchain_id"`
	CTime                  int64  `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func NewCleanBlock() interface{} {
	return &Block{}
}
