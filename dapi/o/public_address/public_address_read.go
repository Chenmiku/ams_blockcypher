package public_address

import ()

func GetByID(id string) (*PublicAddress, error) {
	var w PublicAddress
	return &w, TablePublicAddress.ReadByID(id, &w)
}

func GetByAddress(address string) (*PublicAddress, error) {
	var w PublicAddress
	return &w, TablePublicAddress.ReadOne(map[string]interface{}{
		"address": address,
		"dtime": 0,
	}, &w)
}

func GetAllByWalletID(pageSize int, pageNumber int, sortBy string, sortOrder string, walletId string, pubaddress *[]PublicAddress) (int, error) {
	var where = map[string]interface{}{
		"dtime": 0,
		"wallet_id": walletId,
	}
	exclude := []string{}
	return TablePublicAddress.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, pubaddress, exclude)
}
