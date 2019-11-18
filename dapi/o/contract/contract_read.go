package contract

import ()

func GetByID(id string) (*Contract, error) {
	var w Contract
	return &w, TableContract.ReadByID(id, &w)
}

func GetByAddress(address string) (*Contract, error) {
	var w Contract
	return &w, TableContract.ReadOne(map[string]interface{}{
		"address": address,
		"dtime": 0,
	}, &w)
}
