package private_address

import ()

func GetByID(id string) (*PrivateAddress, error) {
	var w PrivateAddress
	return &w, TablePrivateAddress.ReadByID(id, &w)
}

func GetByAddress(address string) (*PrivateAddress, error) {
	var w PrivateAddress
	return &w, TablePrivateAddress.ReadOne(map[string]interface{}{
		"address": address,
		"dtime": 0,
	}, &w)
}

func GetAllByUserID(pageSize int, pageNumber int, sortBy string, sortOrder string, userId string, priaddress *[]PrivateAddress) (int, error) {
	var where = map[string]interface{}{
		"dtime": 0,
		"user_id": userId,
	}
	exclude := []string{}
	return TablePrivateAddress.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, priaddress, exclude)
}
