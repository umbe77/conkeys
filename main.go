package main

import (
	"conkeys/api"
	"conkeys/config"
	"conkeys/storage"
	"conkeys/storageprovider"
	"conkeys/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()
	fmt.Printf("using %s provider\n", cfg.Provider)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	stg := storageprovider.GetKeyStorage(cfg.Provider)
	usrStorage := storageprovider.GetUserStorage(cfg.Provider)

	adminUser, getUsrErr := usrStorage.Get("admin")
	if getUsrErr != nil {
		adminUser = storage.User{
			UserName: "admin",
			Name:     "Admin",
			LastName: "Admin",
			Email:    "",
		}
		usrStorage.Add(adminUser)
		usrStorage.SetPassword(adminUser.UserName, utility.EncondePassword(cfg.Admin.Password))
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/api/key/*path", api.Get(stg))
	router.GET("/api/keys/*pathSearch", api.GetKeys(stg))
	router.GET("/api/keys", api.GetAllKeys(stg))

	// Create or update key must be an authenticated call
	router.PUT("/api/key/*path", api.Authenticate(), api.Put(stg))
	router.DELETE("/api/key/*path", api.Authenticate(), api.Delete(stg))
	router.GET("/api/checktoken", api.Authenticate(), api.CheckToken())

	router.POST("/api/token", api.Token(usrStorage))

	// TODO: Restrict access to users api only to administrator and not all authenticated users
	router.GET("/api/user/:username", api.Authenticate(), api.GetUser(usrStorage))
	router.GET("/api/users", api.Authenticate(), api.GetUsers(usrStorage))
	router.POST("/api/user", api.Authenticate(), api.AddUser(usrStorage))
	router.PUT("/api/user", api.Authenticate(), api.UpdateUser(usrStorage))
	router.DELETE("/api/user/*username", api.Authenticate(), api.DeleteUser(usrStorage))

	router.Run()
}
