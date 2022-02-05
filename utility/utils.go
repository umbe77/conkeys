package utility

import (
	"crypto/sha512"
	"fmt"
)

func EncondePassword(password string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	encoded := fmt.Sprintf("%x", sha_512.Sum(nil))
	return encoded
}
