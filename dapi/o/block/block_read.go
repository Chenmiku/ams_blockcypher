package block

import ()

func GetByID(id string) (*Block, error) {
	var w Block
	return &w, TableBlock.ReadByID(id, &w)
}

func GetByHash(hash string) (*Block, error) {
	var w Block
	return &w, TableBlock.ReadOne(map[string]interface{}{
		"hash": hash,
		"dtime": 0,
	}, &w)
}

// func GetAllByWalletID(pageSize int, pageNumber int, sortBy string, sortOrder string, walletid string, pubaddress *[]Address) (int, error) {
// 	var where map[string]interface{}
// 	if walletid == "" {
// 		where =  map[string]interface{}{
// 			"dtime": 0,
// 		}
// 	} else {
// 		where = map[string]interface{}{
// 			"dtime": 0,
// 			"wallet_id": walletid,
// 		}
// 	}

// 	exclude := []string{}
// 	return TableAddress.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, pubaddress, exclude)
// }
