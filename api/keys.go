package api

import (
	"conkeys/crypto"
	"conkeys/storage"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Get(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		path := c.Param("path")
		normalizedPath := strings.TrimPrefix(path, "/")
		value, err := stg.Get(normalizedPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if value.T == storage.Crypted {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Cannot access value of %s try to use methods for encrypted values", normalizedPath),
			})
			return
		}
		c.JSON(http.StatusOK, value)
	}
	return f
}

func GetEncrypted(stg storage.KeyStorage, priv *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")
		normalizedPath := strings.TrimPrefix(path, "/")
		value, err := stg.GetEncrypted(normalizedPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if value.T != storage.Crypted {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Try to request not encrypted key",
			})
			return
		}
		encryptedVal := fmt.Sprintf("%s", value.V)
		vBytes, _ := hex.DecodeString(encryptedVal)
		bytes, err := crypto.Decrypt(vBytes, priv)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s", err),
			})
		}
		value.V = string(bytes)
		c.JSON(http.StatusOK, value)
	}
}

func GetKeys(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		pathSearch := c.Param("pathSearch")
		normalizedPathSearch := strings.TrimPrefix(pathSearch, "/")
		res, err := stg.GetKeys(normalizedPathSearch)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	}
	return f
}

func GetAllKeys(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		res := stg.GetAllKeys()
		c.JSON(http.StatusOK, res)
	}
	return f
}

func Put(stg storage.KeyStorage, pub *rsa.PublicKey) gin.HandlerFunc {
	f := func(c *gin.Context) {
		var val storage.Value
		if err := c.ShouldBind(&val); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
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
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": errorMessage,
			})
			return
		}

		path := c.Param("path")
		normalizedPath := strings.TrimPrefix(path, "/")

		if val.T == storage.Crypted {
			encryptedBytes, err := crypto.Encrypt([]byte(fmt.Sprintf("%s", val.V)), pub)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("%s", err),
				})
				return
			}
			encryptedString := hex.EncodeToString(encryptedBytes)
			val.V = "********"
			stg.PutEncrypted(normalizedPath, val, encryptedString)
		} else {
			stg.Put(normalizedPath, val)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("key '%s' saved", normalizedPath),
		})
	}
	return f
}

func Delete(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		path := c.Param("path")
		normalizedPath := strings.TrimPrefix(path, "/")
		stg.Delete(normalizedPath)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("key '%s' removed", normalizedPath),
		})
	}
	return f
}
