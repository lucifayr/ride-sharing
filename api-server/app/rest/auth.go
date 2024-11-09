package rest

import (
	"encoding/base64"
	"ride_sharing_api/app/simulator"
)

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	simulator.S.RandCrypto(b)
	return base64.URLEncoding.EncodeToString(b)
}
