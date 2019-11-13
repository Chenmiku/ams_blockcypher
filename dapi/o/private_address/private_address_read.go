package private_address

func GetByID(id string) (*PrivateAddress, error) {
	var w PrivateAddress
	return &w, TablePrivateAddress.ReadByID(id, &w)
}

func GetByAddress(address string) (*PrivateAddress, error) {
	var w PrivateAddress
	return &w, TablePrivateAddress.ReadOne(map[string]interface{}{
		"address": address,
		"dtime":   0,
	}, &w)
}

func GetAllByWalletID(pageSize int, pageNumber int, sortBy string, sortOrder string, walletid string, priaddress *[]PrivateAddress) (int, error) {
	var where map[string]interface{}
	if walletid == "" {
		where =  map[string]interface{}{
			"dtime": 0,
		}
	} else {
		where = map[string]interface{}{
			"dtime": 0,
			"wallet_id": walletid,
		}
	}
	exclude := []string{}
	return TablePrivateAddress.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, priaddress, exclude)
}
