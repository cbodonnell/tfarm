package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func HMAC(secret []byte, data ...[]byte) string {
	h := hmac.New(sha256.New, secret)
	for _, d := range data {
		h.Write(d)
	}
	signature := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature)
}
