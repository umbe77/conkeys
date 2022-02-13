package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	return privBytes
}

func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)

	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Public Key Not OK")
	}

	return key, nil
}
