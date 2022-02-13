package utility

import (
	"conkeys/crypto"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
)

func EncondePassword(password string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	encoded := fmt.Sprintf("%x", sha_512.Sum(nil))
	return encoded
}

func InitKeyPair(loadKeys func() (*rsa.PrivateKey, *rsa.PublicKey, error), saveKeys func(*rsa.PrivateKey, *rsa.PublicKey) error) (*rsa.PrivateKey, *rsa.PublicKey) {
	var errGenerateKey error
	priv, pub, err := loadKeys()
	if err != nil {
		priv, pub, errGenerateKey = crypto.GenerateKeyPair()
		if errGenerateKey != nil {
			panic(errGenerateKey)
		}
		err = saveKeys(priv, pub)
		if err != nil {
			panic(err)
		}
	}
	return priv, pub
}
