package transaction_output

import ()

func GetByID(id string) (*TransactionOutput, error) {
	var w TransactionOutput
	return &w, TableTransactionOutput.ReadByID(id, &w)
}

// func GetByHash(hash string) (*Transaction, error) {
// 	var w Transaction
// 	return &w, TableTransaction.ReadOne(map[string]interface{}{
// 		"hash": hash,
// 		"dtime": 0,
// 	}, &w)
// }

func GetAllByTransactionID(pageSize int, pageNumber int, sortBy string, sortOrder string, transactionid string, tranOutput *[]TransactionOutput) (int, error) {
	var where map[string]interface{}
	if transactionid == "" {
		where =  map[string]interface{}{
			"dtime": 0,
		}
	} else {
		where = map[string]interface{}{
			"dtime": 0,
			"transaction_id": transactionid,
		}
	}

	exclude := []string{}
	return TableTransactionOutput.ReadPagingSortWithExclude(where, pageSize, pageNumber, sortBy, sortOrder, tranOutput, exclude)
}
