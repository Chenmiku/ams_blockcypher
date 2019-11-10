package wallet

import (
	"ams_system/dapi/api/auth/session"
	"ams_system/dapi/o/wallet"
	"fmt"
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

	s.HandleFunc("/get_all", s.HandleGetAll)
	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
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
	u.Token = "Bearer" + s.Session(w, r).SessionID

	btc := gobcy.API{"Bearer" + s.Session(w, r).SessionID, "bcy", "test"}
	wa := gobcy.Wallet{u.Name, u.Addresses}
	fmt.Println(btc)
	fmt.Println(wa)
	walletName, err := btc.CreateWallet(wa)
	fmt.Printf("Normal Wallet:%+v\n", walletName)
	if err != nil {
		s.SendError(w, err)
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

// Get wallet by id api
func (s *WalletServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetWallet(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete wallet api
func (s *WalletServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetWallet(r)
	if err != nil {
		s.ErrorMessage(w, "wallet_not_found")
		return
	}
	err = wallet.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
