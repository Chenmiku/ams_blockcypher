package private_wallet

import ()

func GetByID(id string) (*PrivateWallet, error) {
	var w PrivateWallet
	return &w, TablePrivateWallet.ReadByID(id, &w)
}

func GetByAddress(address string) (*PrivateWallet, error) {
	var w PrivateWallet
	return &w, TablePrivateWallet.ReadOne(map[string]interface{}{
		"address": address,
		"dtime": 0,
	}, &w)
}

func GetAllByUserID(pageSize int, pageNumber int, sortBy string, sortOrder string, userId string, priWallet *[]PrivateWallet) (int, error) {
	var where = map[string]interface{}{
		"dtime": 0,
		"user_id": userId,
	}
	exclude := []string{}
	return TablePrivateWallet.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, priWallet, exclude)
}
