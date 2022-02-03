package api

import (
	"conkeys/storage"

	"github.com/gin-gonic/gin"
)

func GetUsers(user storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser(user storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func AddUser(user storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateUser(user storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteUser(user storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
