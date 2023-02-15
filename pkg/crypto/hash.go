package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMAC(secret []byte, data ...[]byte) string {
	h := hmac.New(sha256.New, secret)
	for _, d := range data {
		h.Write(d)
	}
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}
