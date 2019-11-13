package transaction

import (
	"ams_system/dapi/o/transaction"
	"http/web"
	"net/http"
	"strconv"
)

type TransactionServer struct {
	web.JsonServer
	*http.ServeMux
}

// create server mux to handle Transaction api
func NewTransactionServer() *TransactionServer {
	var s = &TransactionServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/check_deposit_state", s.HandleCheckDepositState)
	// s.HandleFunc("/get_all", s.HandleGetAll)
	// s.HandleFunc("/get", s.HandleGetByID)
	// s.HandleFunc("/get_by_address", s.HandleGetByAddress)
	// s.HandleFunc("/update", s.HandleUpdateByID)
	// s.HandleFunc("/balance", s.HandleBalance)
	// s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// create transaction api
func (s *TransactionServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var u = &transaction.Transaction{}
	s.MustDecodeBody(r, u)

	// bcy := gobcy.API{"36fd54969a3e499b9bc8f51ee1480d8b", "bcy", "test"}
    // //generate two addresses
    // addr1, err := bcy.GenAddrKeychain()
    // addr2, err := bcy.GenAddrKeychain()
    // //use faucet to fund first
    // _, err = bcy.Faucet(addr1, 3e5)
    // if err != nil {
    //     fmt.Println(err)
    // }
    // //Post New TXSkeleton
    // skel, err := bcy.NewTX(gobcy.TempNewTX(addr1.Address, addr2.Address, 2e5), false)
    // //Sign it locally
    // err = skel.Sign([]string{addr1.Private})
    // if err != nil {
    //     fmt.Println(err)
    // }
    // //Send TXSkeleton
    // skel, err = bcy.SendTX(skel)
    // if err != nil {
    //     fmt.Println(err)
    // }
    // fmt.Printf("%+v\n", skel)

	err := u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, u)
	}
}

// check deposit state api
func (s *TransactionServer) HandleCheckDepositState(w http.ResponseWriter, r *http.Request) {
	// 
}

////get all public address api by walletid
// func (s *PublicAddressServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
// 	walletId := r.URL.Query().Get("wallet_id")
// 	sortBy := r.URL.Query().Get("sort_by")
// 	sortOrder := r.URL.Query().Get("sort_order")

// 	pageSize := StrToInt(r.URL.Query().Get("page_size"))
// 	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

// 	var res = []public_address.PublicAddress{}
// 	count, err := public_address.GetAllByWalletID(pageSize, pageNumber, sortBy, sortOrder, walletId, &res)

// 	if err != nil {
// 		s.SendError(w, err)
// 	} else {
// 		s.SendDataSuccess(w, map[string]interface{}{
// 			"addresss": res,
// 			"count":   count,
// 		})
// 	}
// }

// func (s *PublicAddressServer) mustGetPublicAddress(r *http.Request) (*public_address.PublicAddress, error) {
// 	var id = r.URL.Query().Get("id")
// 	var u, err = public_address.GetByID(id)
// 	if err != nil {
// 		return u, err
// 	} else {
// 		return u, nil
// 	}
// }

// // update public address api
// func (s *PublicAddressServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
// 	var newaddress = &public_address.PublicAddress{}
// 	s.MustDecodeBody(r, newaddress)
// 	u, err := s.mustGetPublicAddress(r)
// 	if err != nil {
// 		s.ErrorMessage(w, "address_not_found")
// 		return
// 	}
// 	err = u.UpdateById(newaddress)
// 	if err != nil {
// 		s.ErrorMessage(w, err.Error())
// 	} else {
// 		result, err := public_address.GetByID(u.ID)
// 		if err != nil {
// 			s.ErrorMessage(w, "address_not_found")
// 			return
// 		}
// 		s.SendDataSuccess(w, result)
// 	}
// }

// // balance public address api
// func (s *PublicAddressServer) HandleBalance(w http.ResponseWriter, r *http.Request) {
// 	address := r.URL.Query().Get("address")
// 	var newaddress = &public_address.PublicAddress{}
// 	s.MustDecodeBody(r, newaddress)
// 	u, err := s.mustGetPublicAddress(r)
// 	if err != nil {
// 		s.ErrorMessage(w, "address_not_found")
// 		return
// 	}

// 	btc := gobcy.API{"36fd54969a3e499b9bc8f51ee1480d8b", "btc", "main"}
// 	addr, err := btc.GetAddrBal(address, nil)
// 	if err != nil {
// 		s.ErrorMessage(w, err.Error())
// 	}

// 	newaddress.TotalRevceived = addr.TotalReceived
// 	newaddress.TotalSent = addr.TotalSent
// 	newaddress.Balance = addr.Balance
// 	newaddress.UnconfirmedBalance = addr.UnconfirmedBalance
// 	newaddress.FinalBalance = addr.FinalBalance
// 	newaddress.ConfirmedTransaction = addr.NumTX
// 	newaddress.UnconfirmedTransaction = addr.UnconfirmedNumTX
// 	newaddress.FinalTransaction = addr.FinalNumTX
// 	err = u.UpdateById(newaddress)
// 	if err != nil {
// 		s.ErrorMessage(w, err.Error())
// 	} else {
// 		result, err := public_address.GetByID(u.ID)
// 		if err != nil {
// 			s.ErrorMessage(w, "address_not_found")
// 			return
// 		}
// 		s.SendDataSuccess(w, result)
// 	}
// }

// // Get public address by address api
// func (s *PublicAddressServer) HandleGetByAddress(w http.ResponseWriter, r *http.Request) {
// 	address := r.URL.Query().Get("address")
// 	u, err := public_address.GetByAddress(address)
// 	if err != nil {
// 		s.ErrorMessage(w, "address_not_found")
// 		return
// 	}
// 	s.SendDataSuccess(w, u)
// }

// // Get public address by id api
// func (s *PublicAddressServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
// 	u, err := s.mustGetPublicAddress(r)
// 	if err != nil {
// 		s.ErrorMessage(w, "address_not_found")
// 		return
// 	}
// 	s.SendDataSuccess(w, u)
// }

// // delete public address api
// func (s *PublicAddressServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
// 	u, err := s.mustGetPublicAddress(r)
// 	if err != nil {
// 		s.ErrorMessage(w, "address_not_found")
// 		return
// 	}
// 	err = public_address.MarkDelete(u.ID)
// 	if err != nil {
// 		s.ErrorMessage(w, err.Error())
// 	} else {
// 		s.Success(w)
// 	}
// }
