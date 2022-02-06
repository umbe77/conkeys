package api

import (
	"conkeys/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Find a standard way to parse a query language (possibly not odata)
func GetUsers(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("userquery")
		users, err := u.GetUsers(query)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("username")
		usr, err := u.Get(userName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, usr)
	}
}

func AddUser(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr storage.User
		if err := c.ShouldBind(&usr); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := u.Add(usr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User Created",
		})
	}
}

func UpdateUser(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr storage.User
		if err := c.ShouldBind(&usr); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := u.Update(usr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User Updated",
		})
	}
}

func DeleteUser(u storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		userName := c.Param("username")

		err := u.Delete(userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User Updated",
		})
	}
}
