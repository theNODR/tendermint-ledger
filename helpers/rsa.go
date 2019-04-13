package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

const PemBlockType = "PUBLIC KEY"

func CreateRSASignatureForBytes(source []byte, sk *rsa.PrivateKey) ([]byte, error) {
	newhash := crypto.SHA256
	pssh:= newhash.New()
	pssh.Write(source)
	hashed := pssh.Sum(nil)

	// var opts rsa.PSSOptions
	// opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example

	// PSS not started (don't understand why) in node-rsa js package
	// result, _ := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, &opts)
	signature, err := rsa.SignPKCS1v15(rand.Reader, sk, newhash, hashed)

	if err != nil {
		return nil, err
	}

	return signature, nil

}

func CreateRSASignatureForString(source string, sk *rsa.PrivateKey) ([]byte, error) {
	return CreateRSASignatureForBytes([]byte(source), sk)
}

func CreateRSASignatureForObject(source interface{}, sk *rsa.PrivateKey) ([]byte, error) {
	bytes, err := json.Marshal(source)

	if err != nil {
		return nil, err
	}

	return CreateRSASignatureForBytes(bytes, sk)
}

func ValidateRSASignatureForBytes(source []byte, base64sign string, pk *rsa.PublicKey) (error) {
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(source)
	digest := pssh.Sum(nil)

	ds, err := base64.StdEncoding.DecodeString(base64sign)

	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(pk, newhash, digest, ds)
	if err != nil {
		return err
	}

	return nil

}

func ValidateRSASignatureForString(source string, base64sign string, pk *rsa.PublicKey) (error) {
	return ValidateRSASignatureForBytes([]byte(source), base64sign, pk)
}

func ValidateRSASignatureForObject(source interface{}, base64sign string, pk *rsa.PublicKey) (error) {
	bytes, err := json.Marshal(source)

	if err != nil {
		return err
	}

	return ValidateRSASignatureForBytes(bytes, base64sign, pk)
}

func DecodeRSAPublicKey(pk string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pk))

	if block == nil {
		return nil, errors.New("invalid public key: not pem encoded")
	}

	if got, want := block.Type, PemBlockType; got != want {
		return nil, errors.New(fmt.Sprintf("invalid public key: invalid format %s", block.Type))
	}

	raw, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	key, ok := raw.(*rsa.PublicKey)

	if !ok {
		return nil, errors.New("invalid parse key result")
	}

	return key, nil
}

func EncodeRSAPublicKeyToBytes(pk crypto.PublicKey) ([]byte, error){
	pubAsn1, err := x509.MarshalPKIXPublicKey(pk)

	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type: PemBlockType,
		Bytes: pubAsn1,
	})

	return pubBytes, nil
}

func EncodeRSAPublicKeyToString(pk crypto.PublicKey) (string, error) {
	pubBytes, err := EncodeRSAPublicKeyToBytes(pk)

	if err != nil {
		return "", err
	}

	return string(pubBytes[:]), nil
}
