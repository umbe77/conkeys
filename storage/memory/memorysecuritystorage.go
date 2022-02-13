package memory

import (
	"crypto/rsa"
	"fmt"
)

type SecurityMemoryStorage struct{}

var cryptoPublicKey *rsa.PublicKey
var cryptoPrivateKey *rsa.PrivateKey

var signinPublicKey *rsa.PublicKey
var signinPrivateKey *rsa.PrivateKey

func (s SecurityMemoryStorage) Init() {}

func (s SecurityMemoryStorage) LoadCryptingPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if cryptoPublicKey == nil || cryptoPrivateKey == nil {
		return nil, nil, fmt.Errorf("No Pair present")
	}
	return cryptoPrivateKey, cryptoPublicKey, nil
}

func (s SecurityMemoryStorage) SaveCryptingPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	cryptoPrivateKey = priv
	cryptoPublicKey = pub
	return nil
}

func (s SecurityMemoryStorage) LoadSigninPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if signinPrivateKey == nil || signinPublicKey == nil {
		return nil, nil, fmt.Errorf("No Pair present")
	}
	return signinPrivateKey, signinPublicKey, nil
}

func (s SecurityMemoryStorage) SaveSigninPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	signinPrivateKey = priv
	signinPublicKey = pub
	return nil
}
