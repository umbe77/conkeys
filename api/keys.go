package api

import (
	"conkeys/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		path := c.Param("path")
		value, err := stg.Get(path)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, value)
	}
	return f
}

func GetKeys(stg storage.KeyStorage) gin.HandlerFunc {
	f := func(c *gin.Context) {
		pathSearch := c.Param("pathSearch")
		res, err := stg.GetKeys(pathSearch)
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
		path := c.Param("path")

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

		stg.Put(path, val)
	}
	return f
}
