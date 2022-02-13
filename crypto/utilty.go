package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func Encrypt(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	cipehrMsg, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return cipehrMsg, nil
}

func Decrypt(msg []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plainMsg, err := rsa.DecryptOAEP(hash, rand.Reader, priv, msg, nil)
	if err != nil {
		return nil, err
	}
	return plainMsg, nil
}
