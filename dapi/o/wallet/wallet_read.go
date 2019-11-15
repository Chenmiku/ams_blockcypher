package wallet

import ()

func GetByID(id string) (*Wallet, error) {
	var w Wallet
	return &w, TableWallet.ReadByID(id, &w)
}

func GetByName(name string) (*Wallet, error) {
	var w Wallet
	return &w, TableWallet.ReadOne(map[string]interface{}{
		"name": name,
		"dtime": 0,
	}, &w)
}

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, userid string, wallet *[]Wallet) (int, error) {
	var where map[string]interface{}
	if userid == "" {
		where =  map[string]interface{}{
			"dtime": 0,
		}
	} else {
		where = map[string]interface{}{
			"dtime": 0,
			"user_id": userid,
		}
	}
	exclude := []string{"token"}
	return TableWallet.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, wallet, exclude)
}
