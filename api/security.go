package api

import (
	"conkeys/storage"
	"crypto/sha512"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var secret = "a;ncieyu89c7y48374q8c ;58yncq8oy5nc58390ycnq835[ncy8*{}@(!&*% &yhn sdUPDA(O*&Y{P(!@*#&BINDUP(O*&(P*)&P($OI@UYO:yu@ad"

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

			token, tkErr := jwt.Parse(authJwt, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signin method")
				}
				return []byte(secret), nil
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

func CheckToken() gin.HandlerFunc {
	f := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
	return f
}

func Token(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr UserLogin
		if err := c.ShouldBind(&usr); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		sha_512 := sha512.New()
		sha_512.Write([]byte(usr.Password))
		pwd := fmt.Sprintf("%x", sha_512.Sum(nil))

		pwd_db, usrErr := u.Get(usr.UserName)
		if usrErr != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User Unauthorized",
			})
			return
		}

		if pwd != pwd_db.Password {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User Unauthorized",
			})
			return
		}

		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer: "UMBE",
			IssuedAt: time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

		tokenString, tkErr := token.SignedString([]byte(secret))
		if tkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})

	}
}
