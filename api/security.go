package api

import (
	"conkeys/storage"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
)

func Authenticate(stg storage.KeyStorage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwt.Payload()
	}
}
