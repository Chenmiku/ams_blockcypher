package transaction

import (
	"ams_system/dapi/o/transaction_input"
	"ams_system/dapi/o/transaction_output"
	"ams_system/dapi/o/address"
	"ams_system/dapi/o/transaction"
	"fmt"
	"http/web"
	"net/http"
	"strconv"

	"github.com/blockcypher/gobcy"
	"ams_system/dapi/config"
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

	s.HandleFunc("/send", s.HandleSend)                                                 
	s.HandleFunc("/deposit", s.HandleDeposit)                                                 
	s.HandleFunc("/with_draw", s.HandleWithDraw)                                                 
	s.HandleFunc("/check_deposit_state", s.HandleCheckDepositState) 
	// s.HandleFunc("/get_all", s.HandleGetAll)
	s.HandleFunc("/get", s.HandleGetByHash)
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// send, receive coin api
func (s *TransactionServer) HandleSend(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	value := StrToInt(r.URL.Query().Get("value"))
	coinType := r.URL.Query().Get("coin_type")
	if strconv.Itoa(value) == "0" || value == 0  || receiver == "" || sender == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}

	sendAddr, err := address.GetByAddress(sender)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}

	_, err = address.GetByAddress(receiver)
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
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	if addr.Balance == 0 || addr.Balance < value {
		s.ErrorMessage(w, "not_enough_fund")
		return
	}

	// Generate address 
	// addr, err := btc.GenAddrKeychain()
	// if err != nil {
	// 	s.ErrorMessage(w, err.Error())
	// 	return
	// }
	// fmt.Println(addr.Address)
	//faucet, only use for test
	// _, err = btc.Faucet(addr, 1000000000000000000)
	// if err != nil {
	// 	s.ErrorMessage(w, err.Error())
	// 	return
	// }
	// fmt.Println("faucet")

	// create new transaction
	skel, err := btc.NewTX(gobcy.TempNewTX(sendAddr.Addr, receiver, value), false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	//sign it
	err = skel.Sign([]string{sendAddr.PrivateKey})
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
	} else {
		s.SendDataSuccess(w, u)
	}
}

// deposit coin api
func (s *TransactionServer) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	privateKey := r.URL.Query().Get("private_key")
	receiver := r.URL.Query().Get("receiver")
	value := StrToInt(r.URL.Query().Get("value"))
	coinType := r.URL.Query().Get("coin_type")
	if strconv.Itoa(value) == "0" || value == 0 || sender == "" || privateKey == "" || receiver == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}
	var sendAddr = &address.Address{}

	_, err := address.GetByAddress(receiver)
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
	// check fund
	addr, err := btc.GetAddrBal(sender, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	if addr.Balance == 0 || addr.Balance < value {
		s.ErrorMessage(w, "not_enough_fund")
		return
	}
	// create new transaction
	skel, err := btc.NewTX(gobcy.TempNewTX(sender, receiver, value), false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	//sign it
	err = skel.Sign([]string{privateKey})
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

	// add sender address to db
	_, err = address.GetByAddress(sender)
	if err != nil {
		sendAddr.Addr = sender
		sendAddr.Balance = addr.Balance
		sendAddr.FinalBalance = addr.FinalBalance
		sendAddr.FinalTransaction = addr.FinalNumTX
		sendAddr.TotalRevceived = addr.TotalReceived
		sendAddr.TotalSent = addr.TotalSent
		sendAddr.UnconfirmedBalance = addr.UnconfirmedBalance
		sendAddr.UnconfirmedTransaction = addr.UnconfirmedNumTX
		sendAddr.ConfirmedTransaction = addr.NumTX
		err = sendAddr.Create()
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
	}

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
	} else {
		s.SendDataSuccess(w, u)
	}
}

// withdraw api
func (s *TransactionServer) HandleWithDraw(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	value := StrToInt(r.URL.Query().Get("value"))
	coinType := r.URL.Query().Get("coin_type")
	if strconv.Itoa(value) == "0" || value == 0 || sender == "" || receiver == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}
	var receiveAddr = &address.Address{}

	sendAddr, err := address.GetByAddress(sender)
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
	// check fund
	addr, err := btc.GetAddrBal(sender, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	if addr.Balance == 0 || addr.Balance < value {
		s.ErrorMessage(w, "not_enough_fund")
		return
	}
	// create new transaction
	skel, err := btc.NewTX(gobcy.TempNewTX(sender, receiver, value), false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	//sign it
	err = skel.Sign([]string{sendAddr.PrivateKey})
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

	// add receiver address to db
	receiverAddr, err := btc.GetAddrBal(receiver, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	_, err = address.GetByAddress(receiver)
	if err != nil {
		receiveAddr.Addr = receiver
		receiveAddr.Balance = receiverAddr.Balance
		receiveAddr.FinalBalance = receiverAddr.FinalBalance
		receiveAddr.FinalTransaction = receiverAddr.FinalNumTX
		receiveAddr.TotalRevceived = receiverAddr.TotalReceived
		receiveAddr.TotalSent = receiverAddr.TotalSent
		receiveAddr.UnconfirmedBalance = receiverAddr.UnconfirmedBalance
		receiveAddr.UnconfirmedTransaction = receiverAddr.UnconfirmedNumTX
		receiveAddr.ConfirmedTransaction = receiverAddr.NumTX
		err = receiveAddr.Create()
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
	}

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
	} else {
		s.SendDataSuccess(w, u)
	}
}

// check deposit state api
func (s *TransactionServer) HandleCheckDepositState(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	coinType := r.URL.Query().Get("coin_type")
	u, err := transaction.GetByHash(hash)
	if err != nil {
		s.ErrorMessage(w, "transaction_not_found")
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
	// check exists transaction
	trans, err := btc.GetTX(hash, nil)
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

		s.SendSuccessMessage(w, "transaction_confirmed", true)
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
