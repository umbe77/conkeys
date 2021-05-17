package api

import (
	"conkeys/config"
	"crypto/sha512"
	"fmt"
	"net/http"
	"strings"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
)

func setUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Cannot access this resource",
	})
}

type AuthPayload struct {
	jwt.Payload
	Usr string `json:"usr,omitempty"`
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		if authHeaderVal, ok := c.Request.Header["Authorization"]; ok {
			authHeader := authHeaderVal[0]

			if !strings.HasPrefix(authHeader, "Bearer ") {
				setUnauthorized(c)
				c.Abort()
				return
			}

			authJwt := strings.Replace(authHeader, "Bearer ", "", 1)

			adminPwd := config.GetConfig().Admin.Password
			sha_512 := sha512.New()
			secret := fmt.Sprintf("%x", sha_512.Sum([]byte(adminPwd)))

			hs := jwt.NewHS512([]byte(secret))
			var pl AuthPayload
			_, err := jwt.Verify([]byte(authJwt), hs, &pl)
			if err != nil {
				setUnauthorized(c)
				c.Abort()
				return
			}

			c.Next()
			return
		}
		setUnauthorized(c)
		c.Abort()
	}
}
