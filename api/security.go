package api

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"conkeys/storage"
	"conkeys/utility"
)

func setUnauthorized(c *fiber.Ctx) error {
	c.Status(fiber.StatusUnauthorized)
	return c.JSON(fiber.Map{
		"error": "Cannot access this resource",
	})
}

type ConkeysClaims struct {
	Adm bool `json:"adm"`
	*jwt.RegisteredClaims
}

type UserLogin struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func Authenticate(priv *rsa.PublicKey, isAdmin bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if authHeader, ok := c.GetReqHeaders()["Authorization"]; ok {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return setUnauthorized(c)
			}

			authJwt := strings.Replace(authHeader, "Bearer ", "", 1)

			token, tkErr := jwt.Parse(authJwt, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signin method")
				}
				return priv, nil
			})

			if tkErr != nil {
				return setUnauthorized(c)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if isAdmin && !claims["adm"].(bool) {
					return setUnauthorized(c)
				}
			}
			return c.Next()
		}

		return setUnauthorized(c)
	}
}

func Token(u storage.UserStorage, pub *rsa.PrivateKey) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var usr UserLogin
		if err := c.BodyParser(&usr); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		pwd := utility.EncondePassword(usr.Password)

		pwd_db, user, usrErr := u.GetPassword(usr.UserName)
		if usrErr != nil {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"error": "User Unauthorized 1",
			})
		}

		if pwd != pwd_db {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"error": "User Unauthorized 2",
			})
		}

		claims := ConkeysClaims{
			user.IsAdmin,
			&jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
				Issuer:    "UMBE",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

		tokenString, tkErr := token.SignedString(pub)
		if tkErr != nil {
			return tkErr
			//TODO: Add Logging
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
			"user":  user,
		})

	}
}
