package wallet

import (
	"ams_system/dapi/api/auth/session"
	"ams_system/dapi/o/wallet"
	"ams_system/dapi/o/private_address"
	"ams_system/dapi/o/public_address"
	"ams_system/dapi/config"
	"http/web"
	"net/http"
	"strconv"

	"github.com/blockcypher/gobcy"
)

type WalletServer struct {
	web.JsonServer
	*http.ServeMux
}

// create server mux to handle user api
func NewWalletServer() *WalletServer {
	var s = &WalletServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.HandleGetAll) //
	s.HandleFunc("/create", s.HandleCreate) //
	s.HandleFunc("/get", s.HandleGetByName) //
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/add_address_to_wallet", s.HandleAddAddress) //
	s.HandleFunc("/remove_address_to_wallet", s.HandleRemoveAddress) //
	s.HandleFunc("/mark_delete", s.HandleMarkDelete) //
	return s
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func (s *WalletServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

// get all wallet api
func (s *WalletServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []wallet.Wallet{}
	count, err := wallet.GetAll(pageSize, pageNumber, sortBy, sortOrder, userId, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"wallets": res,
			"count":   count,
		})
	}
}

// create wallet api
func (s *WalletServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var u = &wallet.Wallet{}
	s.MustDecodeBody(r, u)
	u.Token = config.UserToken

	btc := gobcy.API{u.Token, config.CoinType, config.Chain}
	_, err := btc.CreateWallet(gobcy.Wallet{u.Name, u.Addresses})
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	err = u.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, u)
	}
}

func (s *WalletServer) mustGetWallet(r *http.Request) (*wallet.Wallet, error) {
	var id = r.URL.Query().Get("id")
	var u, err = wallet.GetByID(id)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// update wallet api
func (s *WalletServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newWallet = &wallet.Wallet{}
	s.MustDecodeBody(r, newWallet)
	u, err := s.mustGetWallet(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	err = u.UpdateById(newWallet)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := wallet.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "wallet_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// add address to wallet api
func (s *WalletServer) HandleAddAddress(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	var newWallet = &wallet.Wallet{}
	s.MustDecodeBody(r, newWallet)
	u, err := wallet.GetByName(name)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}

	btc := gobcy.API{u.Token, config.CoinType, config.Chain}
	_, err = btc.AddAddrWallet(name, newWallet.Addresses, false)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	addr := newWallet.Addresses
	for i,_ := range addr {
		var priAdd = &private_address.PrivateAddress{}
		priAdd, err = private_address.GetByAddress(addr[i])
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		priAdd.WalletID = u.ID
		priAdd.WalletName = u.Name
		err = priAdd.UpdateById(priAdd)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		var pubAdd = &public_address.PublicAddress{}
		pubAdd, err = public_address.GetByAddress(addr[i])
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		pubAdd.WalletID = u.ID
		pubAdd.WalletName = u.Name
		err = pubAdd.UpdateById(pubAdd)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
	}

	u.Addresses = append(u.Addresses, newWallet.Addresses...)
	err = u.UpdateById(u)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := wallet.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "wallet_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// remove address to wallet api
func (s *WalletServer) HandleRemoveAddress(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	var newWallet = &wallet.Wallet{}
	s.MustDecodeBody(r, newWallet)
	u, err := wallet.GetByName(name)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}

	btc := gobcy.API{u.Token, config.CoinType, config.Chain}
	err = btc.DeleteAddrWallet(name, newWallet.Addresses)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	addressList := u.Addresses
	removeList := newWallet.Addresses
loop:
	for i := 0; i < len(addressList); i++ {
		addList := addressList[i]
		for _, rem := range removeList {
			if addList == rem {
				addressList = append(addressList[:i], addressList[i+1:]...)
				i--
				continue loop
			}
		}
	}

	for i,_ := range removeList {
		var priAdd = &private_address.PrivateAddress{}
		priAdd, err = private_address.GetByAddress(removeList[i])
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		priAdd.WalletID = ""
		priAdd.WalletName = ""
		err = priAdd.UpdateById(priAdd)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}

		var pubAdd = &public_address.PublicAddress{}
		pubAdd, err = public_address.GetByAddress(removeList[i])
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		pubAdd.WalletID = ""
		pubAdd.WalletName = ""
		err = pubAdd.UpdateById(pubAdd)
		if err != nil {
			s.ErrorMessage(w, err.Error())
			return
		}
	}

	u.Addresses = addressList
	err = u.UpdateById(u)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := wallet.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "wallet_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// Get wallet by name api
func (s *WalletServer) HandleGetByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	u, err := wallet.GetByName(name)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete wallet api
func (s *WalletServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	u, err := wallet.GetByName(name)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}

	btc := gobcy.API{u.Token, config.CoinType, config.Chain}
	err = btc.DeleteWallet(name)
	if err != nil {
		s.ErrorMessage(w, err.Error())
		return
	}

	err = wallet.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
