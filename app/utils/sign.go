package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Usage: Sign("sha1", ...)
// Usage: Sign("sha1=xxxx", ...)
func Sign(alg string, dat []byte, key string) (sig string) {
	sep := "="
	idx := strings.Index(alg, sep)
	if idx >= 0 {
		alg = alg[:idx]
	}

	if alg == "sha1" {
		hash := hmac.New(sha1.New, []byte(key))
		hash.Write(dat)
		sig = hex.EncodeToString(hash.Sum(nil))
	}

	if idx >= 0 && sig != "" {
		return alg + sep + sig
	}

	return
}
