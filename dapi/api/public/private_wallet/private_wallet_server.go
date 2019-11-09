package private_wallet

import (
	"ams_system/dapi/o/private_wallet"
	"http/web"
	"net/http"
	"strconv"

	"github.com/blockcypher/gobcy"
)

type PrivateWalletServer struct {
	web.JsonServer
	*http.ServeMux
}

// create server mux to handle private address api
func NewPrivateWalletServer() *PrivateWalletServer {
	var s = &PrivateWalletServer{
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

// get all wallet api
func (s *PrivateWalletServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []private_wallet.PrivateWallet{}
	count, err := private_wallet.GetAllByUserID(pageSize, pageNumber, sortBy, sortOrder, userId, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"wallets": res,
			"count":   count,
		})
	}
}

// create private key api
func (s *PrivateWalletServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var u = &private_wallet.PrivateWallet{}
	s.MustDecodeBody(r, u)

	btc := gobcy.API{"36fd54969a3e499b9bc8f51ee1480d8b", "btc", "main"}
	addrKeys, err := btc.GenAddrKeychain()
    if err != nil {
		s.ErrorMessage(w, err.Error())
		return
    }

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

func (s *PrivateWalletServer) mustGetPrivateAddress(r *http.Request) (*private_wallet.PrivateWallet, error) {
	var id = r.URL.Query().Get("id")
	var u, err = private_wallet.GetByID(id)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// update private key api
func (s *PrivateWalletServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newWallet = &private_wallet.PrivateWallet{}
	s.MustDecodeBody(r, newWallet)
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	err = u.UpdateById(newWallet)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := private_wallet.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "wallet_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// Get private key by address api
func (s *PrivateWalletServer) HandleGetByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	u, err := private_wallet.GetByAddress(address)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// Get private key by id api
func (s *PrivateWalletServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete private key api
func (s *PrivateWalletServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPrivateAddress(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	err = private_wallet.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
