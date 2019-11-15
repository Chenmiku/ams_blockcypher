package transaction

import (
	"ams_system/dapi/o/transaction_input"
	"ams_system/dapi/o/transaction_output"
	"ams_system/dapi/o/private_address"
	"ams_system/dapi/o/public_address"
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

	s.HandleFunc("/send", s.HandleSend) //                                                
	s.HandleFunc("/deposit", s.HandleDeposit) //                                                
	s.HandleFunc("/check_deposit_state", s.HandleCheckDepositState) //
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
	if strconv.Itoa(value) == "" || value == 0 || sender == "" || receiver == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}
	var sendAddr = &private_address.PrivateAddress{}
	var sendPubAddr = &public_address.PublicAddress{}

	sendAddr, err := private_address.GetByAddress(sender)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	sendPubAddr, err = public_address.GetByAddress(sender)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	if value + 13700 > sendPubAddr.Balance || sendPubAddr.Balance == 0 {
		s.SendError(w, web.ErrBadRequest)
		return
	}
	_, err = private_address.GetByAddress(receiver)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
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

	err = u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, u)
	}
}

// send coin api
func (s *TransactionServer) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get("receiver")
	value := StrToInt(r.URL.Query().Get("value"))
	if strconv.Itoa(value) == "" || value == 0 || receiver == "" {
		s.SendError(w, web.ErrBadRequest)
		return
	}

	var u = &transaction.Transaction{}
	var txInput = &transaction_input.TransactionInput{}
	var txOutput = &transaction_output.TransactionOutput{}
	var txInputs = []transaction_input.TransactionInput{}
	var txOutputs = []transaction_output.TransactionOutput{}
	var sendAddr = &private_address.PrivateAddress{}
	var sendPubAddr = &public_address.PublicAddress{}

	_, err := private_address.GetByAddress(receiver)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	senderAddr, err := btc.GenAddrKeychain()
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
    // use faucet to fund first
	_, err = btc.Faucet(senderAddr, value)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	// create new transaction
	skel, err := btc.NewTX(gobcy.TempNewTX(senderAddr.Address, receiver, value - 13700), false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}
	//sign it
	err = skel.Sign([]string{senderAddr.Private})
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

	// create address on db
	sendAddr.Address = senderAddr.Address
	sendAddr.PublicKey = senderAddr.Public
	sendAddr.PrivateKey = senderAddr.Private
	sendAddr.Wif = senderAddr.Wif
	err = sendAddr.Create() 
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	sendPubAddr.Address = senderAddr.Address
	err = sendPubAddr.Create() 
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
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
	var u = &transaction.Transaction{}
	u, err := transaction.GetByHash(hash)
	if err != nil {
		s.ErrorMessage(w, "transaction_not_found")
		return
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	trans, err := btc.GetTX(hash, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	if trans.Confirmations > 0 {
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
		u.IsCoinBase = false
		err = u.UpdateById(u)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		s.Success(w)
	} else {
		s.ErrorMessage(w, "transaction_not_confirmed")
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