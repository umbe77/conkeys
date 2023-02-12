package api

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"conkeys/crypto"
	"conkeys/storage"
)

func Get(stg storage.KeyStorage) fiber.Handler {
	f := func(c *fiber.Ctx) error {
		path := c.Params("*")
		normalizedPath := strings.TrimPrefix(path, "/")
		value, err := stg.Get(normalizedPath)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if value.T == storage.Crypted {
			value.V = "********"
		}
		c.Status(fiber.StatusOK)
		return c.JSON(value)
	}
	return f
}

func GetEncrypted(stg storage.KeyStorage, priv *rsa.PrivateKey) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Params("*")
		normalizedPath := strings.TrimPrefix(path, "/")
		value, err := stg.GetEncrypted(normalizedPath)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if value.T != storage.Crypted {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": "Try to request not encrypted key",
			})
		}
		encryptedVal := fmt.Sprintf("%s", value.V)
		vBytes, _ := hex.DecodeString(encryptedVal)
		bytes, err := crypto.Decrypt(vBytes, priv)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": fmt.Sprintf("%s", err),
			})
		}
		value.V = string(bytes)
		return c.JSON(value)
	}
}

func GetKeys(stg storage.KeyStorage) fiber.Handler {
	f := func(c *fiber.Ctx) error {
		pathSearch := c.Params("*")
		normalizedPathSearch := strings.TrimPrefix(pathSearch, "/")
		res, err := stg.GetKeys(normalizedPathSearch)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(res)
	}
	return f
}

func GetAllKeys(stg storage.KeyStorage) fiber.Handler {
	f := func(c *fiber.Ctx) error {
		res, err := stg.GetAllKeys()
		if err != nil {
			//TODO: Add logging
			return err
		}
		return c.JSON(res)
	}
	return f
}

func Put(stg storage.KeyStorage, pub *rsa.PublicKey) fiber.Handler {
	f := func(c *fiber.Ctx) error {
		var val storage.Value
		if err := c.BodyParser(&val); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if !val.CheckType() {
			var errorMessage string
			valT, err := val.T.ToString()
			if err != nil {
				errorMessage = err.Error()
			} else {
				errorMessage = fmt.Sprintf("%v is not of type %v", val.V, valT)
			}
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		path := c.Params("*")
		normalizedPath := strings.TrimPrefix(path, "/")

		if val.T == storage.Crypted {
			encryptedBytes, err := crypto.Encrypt([]byte(fmt.Sprintf("%s", val.V)), pub)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"error": fmt.Sprintf("%s", err),
				})
			}
			encryptedString := hex.EncodeToString(encryptedBytes)
			val.V = "********"
			stg.PutEncrypted(normalizedPath, val, encryptedString)
		} else {
			stg.Put(normalizedPath, val)
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("key '%s' saved", normalizedPath),
		})
	}
	return f
}

func Delete(stg storage.KeyStorage) fiber.Handler {
	f := func(c *fiber.Ctx) error {
		path := c.Params("*")
		normalizedPath := strings.TrimPrefix(path, "/")
		err := stg.Delete(normalizedPath)
		if err != nil {
			return err
			//TODO: Add Logging
		}
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("key '%s' removed", normalizedPath),
		})
	}
	return f
}
