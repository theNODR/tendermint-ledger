package helpers

import (
	"crypto/md5"
	"crypto/elliptic"
)

var hasher = md5.New
var curver = elliptic.P256
var ecdsaKeySize = 32
