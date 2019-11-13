package public

import (
	"ams_system/dapi/api/public/org"
	"ams_system/dapi/api/public/private_address"
	"ams_system/dapi/api/public/public_address"
	"ams_system/dapi/api/public/wallet"
	"ams_system/dapi/api/public/transaction"
	"ams_system/dapi/config"
	"http/web"
	"net/http"
)

type PublicServer struct {
	web.JsonServer
	*http.ServeMux
}

func NewPublicServer(pc *config.ProjectConfig) *PublicServer {
	var s = &PublicServer{
		ServeMux: http.NewServeMux(),
	}

	s.Handle("/org/", http.StripPrefix("/org", org.NewOrgServer()))
	s.Handle("/private_address/", http.StripPrefix("/private_address", private_address.NewPrivateAddressServer()))
	s.Handle("/public_address/", http.StripPrefix("/public_address", public_address.NewPublicAddressServer()))
	s.Handle("/wallet/", http.StripPrefix("/wallet", wallet.NewWalletServer()))
	s.Handle("/transaction/", http.StripPrefix("/transaction", transaction.NewTransactionServer()))
	return s
}
