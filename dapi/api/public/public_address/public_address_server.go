package public_address

import (
	"ams_system/dapi/o/public_address"
	"http/web"
	"net/http"
	"strconv"
)

type PublicAddressServer struct {
	web.JsonServer
	*http.ServeMux
}

// create server mux to handle public address api
func NewPublicAddressServer() *PublicAddressServer {
	var s = &PublicAddressServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/get_all", s.HandleGetAll)
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

//get all public address api by walletid
func (s *PublicAddressServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	walletId := r.URL.Query().Get("wallet_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []public_address.PublicAddress{}
	count, err := public_address.GetAllByWalletID(pageSize, pageNumber, sortBy, sortOrder, walletId, &res)

	if err != nil {
		s.SendError(w, err)
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"addresss": res,
			"count":   count,
		})
	}
}

func (s *PublicAddressServer) mustGetPublicAddress(r *http.Request) (*public_address.PublicAddress, error) {
	var id = r.URL.Query().Get("id")
	var u, err = public_address.GetByID(id)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// update public address api
func (s *PublicAddressServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newaddress = &public_address.PublicAddress{}
	s.MustDecodeBody(r, newaddress)
	u, err := s.mustGetPublicAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = u.UpdateById(newaddress)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := public_address.GetByID(u.ID)
		if err != nil {
			s.ErrorMessage(w, "address_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

// Get public address by address api
func (s *PublicAddressServer) HandleGetByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	u, err := public_address.GetByAddress(address)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// Get public address by id api
func (s *PublicAddressServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPublicAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	s.SendDataSuccess(w, u)
}

// delete public address api
func (s *PublicAddressServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	u, err := s.mustGetPublicAddress(r)
	if err != nil {
		s.ErrorMessage(w, "address_not_found")
		return
	}
	err = public_address.MarkDelete(u.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
