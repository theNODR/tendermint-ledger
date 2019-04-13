package helpers

import (
	"crypto/sha256"
	"encoding/base64"
)

func CreateHash(source []byte) string {
	hasher := sha256.New224()
	hasher.Write(source)
	shaBytes := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return shaBytes
}
