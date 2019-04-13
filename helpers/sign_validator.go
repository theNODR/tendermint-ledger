package helpers

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"encoding/base64"
	"math/big"
	"hash"
	"encoding/json"
)

type Validator interface {
	ValidateBytes(source []byte, base64sign string) error
	ValidateString(source string, base64sign string) error
	ValidateObject(source interface{}, base64sign string) error
}

type ecdsaValidator struct {
	publicKey	*ecdsa.PublicKey
}

func NewECSDAValidator(pk string) (*ecdsaValidator, error) {
	blockPub, _ := pem.Decode([]byte(pk))
	if blockPub == nil {
		return nil, errors.New("invalid encoded public key")
	}
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}

	publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, errors.New("invalid key format")
	}

	return &ecdsaValidator{
		publicKey: publicKey,
	}, nil
}

func (v *ecdsaValidator) ValidateObject(source interface{}, base64sign string) error {
	bytes, err := json.Marshal(source)

	if err != nil {
		return err
	}

	return v.ValidateBytes(bytes, base64sign)
}

func (v *ecdsaValidator) ValidateString(source string, base64sign string) error {
	return v.ValidateBytes([]byte(source), base64sign)
}

func (v *ecdsaValidator) ValidateBytes(source []byte, base64sign string) error {
	ds, err := base64.StdEncoding.DecodeString(base64sign)
	if err != nil {
		return err
	}

	if len(ds) != 2 * ecdsaKeySize {
		return errors.New("invalid ECDSA signature size")
	}

	r := big.NewInt(0).SetBytes(ds[:ecdsaKeySize])
	s := big.NewInt(0).SetBytes(ds[ecdsaKeySize:])

	var h hash.Hash
	h = hasher()
	h.Write(source)

	if verifyStatus := ecdsa.Verify(v.publicKey, h.Sum(nil), r, s); verifyStatus == true {
		return nil
	} else {
		return errors.New("invalid ECDSA signature")
	}
}
