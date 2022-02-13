package api

import (
	"conkeys/storage"
	"conkeys/utility"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func setUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Cannot access this resource",
	})
}

// type AuthPayload struct {
// 	jwt.Payload
// 	Usr string `json:"usr,omitempty"`
// }

type UserLogin struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func Authenticate(priv *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {

		if authHeaderVal, ok := c.Request.Header["Authorization"]; ok {
			authHeader := authHeaderVal[0]
			if !strings.HasPrefix(authHeader, "Bearer ") {
				setUnauthorized(c)
				c.Abort()
				return
			}

			authJwt := strings.Replace(authHeader, "Bearer ", "", 1)

			token, tkErr := jwt.Parse(authJwt, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signin method")
				}
				return priv, nil
			})

			if tkErr != nil {
				setUnauthorized(c)
				c.Abort()
				return
			}

			if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Next()
				return
			}
		}
		setUnauthorized(c)
		c.Abort()
	}
}

func Token(u storage.UserStorage, pub *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr UserLogin
		if err := c.ShouldBind(&usr); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		pwd := utility.EncondePassword(usr.Password)

		pwd_db, usrErr := u.GetPassword(usr.UserName)
		if usrErr != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User Unauthorized",
			})
			return
		}

		if pwd != pwd_db {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User Unauthorized",
			})
			return
		}

		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "UMBE",
			IssuedAt:  time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

		tokenString, tkErr := token.SignedString(pub)
		if tkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})

	}
}
