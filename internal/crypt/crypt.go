package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/HelliWrold1/cloud/internal/config"
)

func SetPwd(password string) string {
	mac := hmac.New(sha256.New, []byte(config.Get().Crypt.Sha256Key))
	mac.Write([]byte(password))
	sum := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}

func CheckPwd(password string, hashedPwd string) bool {
	mac := hmac.New(sha256.New, []byte(config.Get().Crypt.Sha256Key))
	mac.Write([]byte(password))
	sum := mac.Sum(nil)
	base64.URLEncoding.EncodeToString(sum)
	if SetPwd(password) == hashedPwd {
		return true
	}
	return false
}
