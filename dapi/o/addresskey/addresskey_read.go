package addresskey

import (
)

func GetByAddress(address string) (*AddressKey, error) {
	var w AddressKey
	return &w, TableAddressKey.ReadOne(map[string]interface{}{
		"addr": address,
		"dtime": 0,
	}, &w)
}

