package helpers

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"hash"
	"math/big"
	"io"
	"encoding/json"
)

type KeyPair interface {
	SignObject(data interface{}) ([]byte, error)
	SignBytes(data []byte) ([]byte, error)
	SignString(data string) ([]byte, error)
	PublicKey() (string, error)
	PrivateKey() (string, error)
}

type ecdsaKeyPair struct {
	privateKey	*ecdsa.PrivateKey
	publicKey	*ecdsa.PublicKey

	pk			*string
	sk			*string
}

func NewECSDAKeyPair() (KeyPair, error) {
	publicCurve := curver()

	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(publicCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ecdsaKeyPair{
		privateKey: privateKey,
		publicKey: &privateKey.PublicKey,
	}, nil
}

func (e *ecdsaKeyPair) SignObject(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return e.SignBytes(bytes)
}

func (e *ecdsaKeyPair) SignBytes(data []byte) ([]byte, error) {
	return e.SignString(string(data[:]))
}

func (e *ecdsaKeyPair) SignString(data string) ([]byte, error) {
	var h hash.Hash
	h = hasher()
	r := big.NewInt(0)
	s := big.NewInt(0)

	io.WriteString(h, data)
	signhash := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, e.privateKey, signhash)

	if err != nil {
		return nil, err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	return signature, nil
}

func (e *ecdsaKeyPair) PublicKey() (string, error) {
	if e.pk != nil {
		return *e.pk, nil
	}

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return "", err
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	s := string(pemEncodedPub[:])
	e.pk = &s
	return s, nil
}

func (e *ecdsaKeyPair) PrivateKey() (string, error) {
	if e.sk !=  nil {
		return *e.sk, nil
	}

	x509Encoded, err := x509.MarshalECPrivateKey(e.privateKey)
	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	s := string(pemEncoded[:])
	e.sk = &s
	return s, nil
}
