package transaction

import ()

func GetByID(id string) (*Transaction, error) {
	var w Transaction
	return &w, TableTransaction.ReadByID(id, &w)
}

func GetByHash(hash string) (*Transaction, error) {
	var w Transaction
	return &w, TableTransaction.ReadOne(map[string]interface{}{
		"hash": hash,
		"dtime": 0,
	}, &w)
}
