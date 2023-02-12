package storage

import "crypto/rsa"

type SecurityStorage interface {
	LoadCryptingPair() (*rsa.PrivateKey, *rsa.PublicKey, error)
	SaveCryptingPair(*rsa.PrivateKey, *rsa.PublicKey) error
	LoadSigninPair() (*rsa.PrivateKey, *rsa.PublicKey, error)
	SaveSigninPair(*rsa.PrivateKey, *rsa.PublicKey) error
}
