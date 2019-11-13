package private_address

import (
	"ams_system/dapi/o/private_address"
	"ams_system/dapi/o/wallet"
	"http/web"
	"net/http"
	"strconv"

	"github.com/blockcypher/gobcy"
)

type PrivateAddressServer struct {
	web.JsonServer
	*http.ServeMux
}

type WalletAddress struct {
	Wallet wallet.Wallet
	Adress private_address.PrivateAddress
}

// create server mux to handle private address api
func NewPrivateAddressServer() *PrivateAddressServer {
	var s = &PrivateAddressServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.HandleGetAll)
	s.HandleFunc("/generate", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/get_by_address", s.HandleGetByAddress)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// get all address api by userid
func (s *PrivateAddressServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	walletId := r.URL.Query().Get("wallet_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []private_address.PrivateAddress{}
	count, err := private_address.GetAllByWalletID(pageSize, pageNumber, sortBy, sortOrder, walletId, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"addresss": res,
			"count":    count,
		})
	}
}

// create private address api
func (s *PrivateAddressServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	walletid := r.URL.Query().Get("wallet_id")
	var wa = &wallet.Wallet{}
	var err error
	if walletid != "" {
		wa, err = wallet.GetByID(walletid)
	}

	var u = &private_address.PrivateAddress{}
	s.MustDecodeBody(r, u)

	btc := gobcy.API{"36fd54969a3e499b9bc8f51ee1480d8b", "btc", "main"}
	addrKeys, err := btc.GenAddrKeychain()
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	wa.Addresses = append(wa.Addresses, addrKeys.Address)
	err = wa.UpdateById(wa)

	u.WalletID = walletid
	u.Address = addrKeys.Address
	u.PublicKey = addrKeys.Public
	u.PrivateKey = addrKeys.Private
	u.Wif = addrKeys.Wif
	err = u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, u)
	}
}

func (s *PrivateAddressServer) mustGetPrivateAddress(r *http.Request) (*private_address.PrivateAddress, error) {
	var id = r.URL.Query().Get("id")
	var u, err = private_address.GetByID(id)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// update private address api
func (s *PrivateAddressServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newaddress = &private_address.PrivateAddress{}
	s.MustDecodeBody(r, newaddress)
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = u.UpdateById(newaddress)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := private_address.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// Get private address by address api
func (s *PrivateAddressServer) HandleGetByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	u, err := private_address.GetByAddress(address)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// Get private address by id api
func (s *PrivateAddressServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete private address api
func (s *PrivateAddressServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = private_address.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
