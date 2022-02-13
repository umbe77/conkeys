package api

import (
	"conkeys/storage"
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

// TODO: Create a new method to get Crypto typed keys

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

func Put(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		var val storage.Value
		if err := c.ShouldBind(&val); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
		}
		// TODO: check if not a Crypted value and crypt the value
		path := c.Param("path")
		normalizedPath := strings.TrimPrefix(path, "/")

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

		stg.Put(normalizedPath, val)
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
