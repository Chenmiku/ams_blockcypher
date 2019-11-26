package transaction

import (
	"ams_system/dapi/o/transaction_input"
	"ams_system/dapi/o/transaction_output"
	"ams_system/dapi/o/address"
	"ams_system/dapi/o/addresskey"
	"ams_system/dapi/o/transaction"
	"fmt"
	"http/web"
	"net/http"
	"strconv"
	"time"

	"github.com/blockcypher/gobcy"
	"ams_system/dapi/config"
)

type TransactionServer struct {
	web.JsonServer
	*http.ServeMux
}

type DepositStateResult struct {
	CoinType  string `json:"coin_type"`
	CoinValue int `json:"coin_value"`
	Confirm	  bool 	`json:"confirm"`
	Message	  string 	`json:"message"`
}

// create server mux to handle Transaction api
func NewTransactionServer() *TransactionServer {
	var s = &TransactionServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/send_to_polebit", s.HandleSend)                                                                                                  
	s.HandleFunc("/check_deposit_state", s.HandleCheckDepositState) 
	s.HandleFunc("/check_deposit_state_by_address", s.HandleCheckDepositStateByAddress) 
	s.HandleFunc("/get", s.HandleGetByHash)
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// send coin to polebit wallet api
func (s *TransactionServer) HandleSend(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	value := StrToInt(r.URL.Query().Get("value"))
	coinType := r.URL.Query().Get("coin_type")
	fmt.Println(value)
	if value == 0  || receiver == "" || sender == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}

	ad, err := address.GetByAddress(sender)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	addrKey, err := addresskey.GetByAddress(sender)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}

	switch coinType {
	case "btc":
		config.CoinType = "btc"
	case "eth":
		config.CoinType = "eth"
	case "":
		config.CoinType = "btc"
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	//check fund
	addr, err := btc.GetAddrBal(sender, nil)
	if addr.Balance == 0 || addr.Balance < value {
		s.ErrorMessage(w, "not_enough_fund")
		return
	}
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	// create new transaction
	skel, err := btc.NewTX(gobcy.TempNewTX(addr.Address, receiver, value), false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	//sign it
	err = skel.Sign([]string{addrKey.PrivateKey})
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	// send transaction
	skel, err = btc.SendTX(skel)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	fmt.Println("success")
	fmt.Println(skel.Trans)

	// create txoutput on db
	txO := skel.Trans.Outputs
	for i,_ := range txO {
		txOutput.Value = txO[i].Value
		txOutput.ScriptType = txO[i].ScriptType
		txOutput.Script = txO[i].Script
		txOutput.Addresses = txO[i].Addresses
		err = txOutput.Create()
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
		txOutputs = append(txOutputs, *txOutput)
	}
	//create txinput on db 
	txI := skel.Trans.Inputs
	for i,_ := range txI {
		txInput.PreviousHash = txI[i].PrevHash
		txInput.OutputIndex = txI[i].OutputIndex
		txInput.OutputValue = txI[i].OutputValue
		txInput.ScriptType = txI[i].ScriptType
		txInput.Script = txI[i].Script
		txInput.Addresses = txI[i].Addresses
		err = txInput.Create()
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
		txInputs = append(txInputs, *txInput)
	}
	// create tx on db 
	u.Hash = skel.Trans.Hash
	u.BlockHeight = skel.Trans.BlockHeight
	u.TotalBlock = skel.Trans.Confirmations
	u.TotalExchanged = skel.Trans.Total
	u.Fees = skel.Trans.Fees
	u.Size = skel.Trans.Size
	u.DoubleSpend = skel.Trans.DoubleSpend
	u.Inputs = txInputs
	u.Outputs = txOutputs
	u.Addresses = skel.Trans.Addresses
	u.ToSign = skel.ToSign
	u.Signatures = skel.Signatures
	u.PublicKeys = skel.PubKeys
	// if config.CoinType == "eth" {
	// 	u.GasUsed = skel.Trans.GasUsed
	// 	u.GasPrice = skel.Trans.GasPrice
	// 	u.GasLimit = skel.Trans.GasLimit
	// }
	err = u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} 
	// else {
	// 	s.SendDataSuccess(w, u)
	// }

	// auto check transaction

	// ticker := time.NewTicker(3 * time.Second)
	// defer ticker.Stop()
    // done := make(chan bool)
    
    // go func ()  {
    //     time.Sleep(60 * time.Second)
    //     done <- true
    // }()

	// for {
	// 	select {
    //     case _ = <-ticker.C:
    //         trans, err := btc.GetTX(u.Hash, nil)
	// 		if err != nil {
	// 			s.ErrorMessage(w, "transaction_not_found")
	// 			return
	// 		}

	// 		if trans.Confirmations > 0 {
	// 			// update transaction on db
	// 			u.Hash = trans.Hash
	// 			u.BlockHeight = trans.BlockHeight
	// 			u.BlockHash = trans.BlockHash
	// 			u.TotalBlock = trans.Confirmations
	// 			u.TotalExchanged = trans.Total
	// 			u.Fees = trans.Fees
	// 			u.Size = trans.Size
	// 			u.Version = trans.Ver
	// 			u.DoubleSpend = trans.DoubleSpend
	// 			u.ConfirmedTime = trans.Confirmed.Unix()
	// 			u.InputsTransaction = trans.VinSize
	// 			u.OutputsTransaction = trans.VoutSize
	// 			u.Addresses = trans.Addresses
	// 			// if config.CoinType == "eth" {
	// 			// 	u.GasUsed = trans.GasUsed
	// 			// 	u.GasPrice = trans.GasPrice
	// 			// 	u.GasLimit = trans.GasLimit
	// 			// }
	// 			err = u.UpdateById(u)
	// 			if err != nil {
	// 				s.ErrorMessage(w, err.Error())
	// 				return
	// 			}
	// 			s.SendSuccessMessage(w, "transaction_confirmed", true)
	// 		} 
	// 		// else {
	// 		// 	s.SendSuccessMessage(w, "transaction_not_confirmed", false)
	// 		// }
    //     case <-done:
	// 		fmt.Println("Done!")
	// 		return
	// 	}
	// }
	
	time.Sleep(30 * time.Second)
	trans, err := btc.GetTX(u.Hash, nil)
	if err != nil {
		s.ErrorMessage(w, "transaction_not_found")
		return
	}

	if trans.Confirmations > 0 {
		// update transaction on db
		u.Hash = trans.Hash
		u.BlockHeight = trans.BlockHeight
		u.BlockHash = trans.BlockHash
		u.TotalBlock = trans.Confirmations
		u.TotalExchanged = trans.Total
		u.Fees = trans.Fees
		u.Size = trans.Size
		u.Version = trans.Ver
		u.DoubleSpend = trans.DoubleSpend
		u.ConfirmedTime = trans.Confirmed.Unix()
		u.InputsTransaction = trans.VinSize
		u.OutputsTransaction = trans.VoutSize
		u.Addresses = trans.Addresses
		// if config.CoinType == "eth" {
		// 	u.GasUsed = trans.GasUsed
		// 	u.GasPrice = trans.GasPrice
		// 	u.GasLimit = trans.GasLimit
		// }
		err = u.UpdateById(u)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		// balance
		time.Sleep(30 * time.Second)
		addre, err := btc.GetAddrBal(sender, nil)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		ad.TotalRevceived = addre.TotalReceived
		ad.TotalSent = addre.TotalSent
		ad.Balance = addre.Balance
		ad.UnconfirmedBalance = addre.UnconfirmedBalance
		ad.FinalBalance = addre.FinalBalance
		ad.ConfirmedTransaction = addre.NumTX
		ad.UnconfirmedTransaction = addre.UnconfirmedNumTX
		ad.FinalTransaction = addre.FinalNumTX
		err = ad.UpdateById(ad)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		s.SendSuccessMessage(w, "transaction_confirmed", true)
	} else {
		s.SendSuccessMessage(w, "transaction_not_confirmed", false)
	}
}

// check deposit state api
func (s *TransactionServer) HandleCheckDepositState(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	coinType := r.URL.Query().Get("coin_type")

	switch coinType {
	case "btc":
		config.CoinType = "btc"
	case "eth":
		config.CoinType = "eth"
	case "":
		config.CoinType = "btc"
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	// check confirm transaction every 3 second
	//ticker := time.NewTicker(3 * time.Second)

	time.Sleep(30 * time.Second)
	trans, err := btc.GetTX(hash, nil)
	if err != nil {
		s.ErrorMessage(w, "transaction_not_found")
		return
	}

	if trans.Confirmations > 0 {
		// save transaction on db
		u, err := transaction.GetByHash(hash)
		if err == nil {
			u.Hash = trans.Hash
			u.BlockHeight = trans.BlockHeight
			u.BlockHash = trans.BlockHash
			u.TotalBlock = trans.Confirmations
			u.TotalExchanged = trans.Total
			u.Fees = trans.Fees
			u.Size = trans.Size
			u.Version = trans.Ver
			u.DoubleSpend = trans.DoubleSpend
			u.ConfirmedTime = trans.Confirmed.Unix()
			u.InputsTransaction = trans.VinSize
			u.OutputsTransaction = trans.VoutSize
			u.Addresses = trans.Addresses
			// if config.CoinType == "eth" {
			// 	u.GasUsed = trans.GasUsed
			// 	u.GasPrice = trans.GasPrice
			// 	u.GasLimit = trans.GasLimit
			// }
			err = u.UpdateById(u)
			if err != nil {
				s.ErrorMessage(w, err.Error())
				return
			}
		}
		s.SendSuccessMessage(w, "transaction_confirmed", true)
	} else {
		s.SendSuccessMessage(w, "transaction_not_confirmed", false)
	}
}

// check deposit state by address api
func (s *TransactionServer) HandleCheckDepositStateByAddress(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	coinType := r.URL.Query().Get("coin_type")
	confirm := false
	coinValue := 0
	ad, err := address.GetByAddress(addr)
	if err != nil { 
		s.ErrorMessage(w, err.Error())
		return
	}

	// check coin type
	switch coinType {
	case "btc":
		config.CoinType = "btc"
	case "eth":
		config.CoinType = "eth"
	case "":
		config.CoinType = "btc"
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	time.Sleep(30 * time.Second)
	addre, err := btc.GetAddrBal(addr, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	if ad.Balance != addre.Balance {
		confirm = true
		coinValue = addre.Balance - ad.Balance
	}

	ad.TotalRevceived = addre.TotalReceived
	ad.TotalSent = addre.TotalSent
	ad.Balance = addre.Balance
	ad.UnconfirmedBalance = addre.UnconfirmedBalance
	ad.FinalBalance = addre.FinalBalance
	ad.ConfirmedTransaction = addre.NumTX
	ad.UnconfirmedTransaction = addre.UnconfirmedNumTX
	ad.FinalTransaction = addre.FinalNumTX
	err = ad.UpdateById(ad)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	if confirm == true {
		result := &DepositStateResult{}
		result.CoinType = config.CoinType
		result.CoinValue = coinValue
		result.Confirm = confirm
		result.Message = "transaction_confirmed"

		s.SendDataSuccess(w, result)
	} else {
		s.SendSuccessMessage(w, "transaction_not_confirmed", false)
	}
}

// Get transaction by hash api
func (s *TransactionServer) HandleGetByHash(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	u, err := transaction.GetByHash(hash)
	if err != nil {
		s.ErrorMessage(w, "transaction_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}
