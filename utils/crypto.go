package utils

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

var ErrWrongKeyAlgo = errors.New("wrong key algorithm")
var ErrDecodePubKey = errors.New("failed to decode Travis public key")

type ErrInvalidKeyType struct {
	Type string
}

func (e *ErrInvalidKeyType) Error() string {
	return fmt.Sprintf("invalid key type: %s", e.Type)
}

func TravisParsePubKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, ErrDecodePubKey
	} else if block.Type != "PUBLIC KEY" {
		return nil, &ErrInvalidKeyType{Type: block.Type}
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	k, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrWrongKeyAlgo
	}
	return k, nil
}

func DataDigest(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	return hash.Sum(nil)
}