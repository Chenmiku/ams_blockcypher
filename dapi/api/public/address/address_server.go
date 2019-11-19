package address

import (
	"ams_system/dapi/o/address"
	"http/web"
	"net/http"
	"strconv"
	"github.com/blockcypher/gobcy"
	"ams_system/dapi/config"
)

type AddressServer struct {
	web.JsonServer
	*http.ServeMux
}

// create server mux to handle public address api
func NewAddressServer() *AddressServer {
	var s = &AddressServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.HandleGetAll) 
	s.HandleFunc("/generate", s.HandleCreate) 
	s.HandleFunc("/get_by_address", s.HandleGetByAddress) 
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/balance", s.HandleBalance) 
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
} 

//get all public address api by user
func (s *AddressServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	userid := StrToInt(r.URL.Query().Get("user_id"))
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []address.Address{}
	count, err := address.GetAllByUser(pageSize, pageNumber, sortBy, sortOrder, userid, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"addresses": res,
			"count":   count,
		})
	}
}

// generate address api
func(s *AddressServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userid := StrToInt(r.URL.Query().Get("user_id"))
	coinType := r.URL.Query().Get("coin_type")
	var err error

	var u = &address.Address{}

	switch coinType {
	case "btc":
		config.CoinType = "btc"
	case "eth":
		config.CoinType = "eth"
	case "":
		config.CoinType = "btc"
	}

	btc := gobcy.API{config.UserToken, config.CoinType, config.Chain}
	addrKeys, err := btc.GenAddrKeychain()
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	u.Addr = addrKeys.Address
	u.PublicKey = addrKeys.Public
	u.PrivateKey = addrKeys.Private
	u.Wif = addrKeys.Wif
	u.CoinType = config.CoinType
	u.UserID = userid
	err = u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, u)
	}
}

func (s *AddressServer) mustGetAddress(r *http.Request) (*address.Address, error) {
	var id = r.URL.Query().Get("id")
	var u, err = address.GetByID(id)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// update address api
func (s *AddressServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newaddress = &address.Address{}
	s.MustDecodeBody(r, newaddress)
	u, err := s.mustGetAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = u.UpdateById(newaddress)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := address.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// balance address api
func (s *AddressServer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	coinType := r.URL.Query().Get("coin_type")
	var priAddress = &address.Address{}

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
	addre, err := btc.GetAddrBal(addr, nil)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	ad, err := address.GetByAddress(addr)
	ad.CoinType = config.CoinType
	ad.TotalRevceived = addre.TotalReceived
	ad.TotalSent = addre.TotalSent
	ad.Balance = addre.Balance
	ad.UnconfirmedBalance = addre.UnconfirmedBalance
	ad.FinalBalance = addre.FinalBalance
	ad.ConfirmedTransaction = addre.NumTX
	ad.UnconfirmedTransaction = addre.UnconfirmedNumTX
	ad.FinalTransaction = addre.FinalNumTX
	if err != nil {
		// create address balance
		priAddress.Addr = addr
		priAddress.CoinType = config.CoinType
		priAddress.TotalRevceived = addre.TotalReceived
		priAddress.TotalSent = addre.TotalSent
		priAddress.Balance = addre.Balance
		priAddress.UnconfirmedBalance = addre.UnconfirmedBalance
		priAddress.FinalBalance = addre.FinalBalance
		priAddress.ConfirmedTransaction = addre.NumTX
		priAddress.UnconfirmedTransaction = addre.UnconfirmedNumTX
		priAddress.FinalTransaction = addre.FinalNumTX
		err = priAddress.Create()
		if err != nil {
			s.ErrorMessage(w, err.Error())
		} else {
			s.SendDataSuccess(w, priAddress)
		}
	} else {
		// update public address balance
		err = ad.UpdateById(ad)
		if err != nil {
			s.ErrorMessage(w, err.Error())
		} else {
			s.SendDataSuccess(w, ad)
		}
	}
}

// Get address by address api
func (s *AddressServer) HandleGetByAddress(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	u, err := address.GetByAddress(addr)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete address api
func (s *AddressServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = address.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
