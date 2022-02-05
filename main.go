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
	fmt.Println(cfg.Provider)
	r := gin.Default()
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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/key/*path", api.Get(stg))
	r.GET("/api/keys/*pathSearch", api.GetKeys(stg))
	r.GET("/api/keys", api.GetAllKeys(stg))

	// Create or update key must be an authenticated call
	r.PUT("/api/key/*path", api.Authenticate(), api.Put(stg))
	r.DELETE("/api/key/*path", api.Authenticate(), api.Delete(stg))
	r.GET("/api/checktoken", api.Authenticate(), api.CheckToken())

	r.POST("/api/token", api.Token(usrStorage))

	r.GET("/api/user/*username", api.Authenticate(), api.GetUser(usrStorage))
	r.GET("/api/users/*userquery", api.Authenticate(), api.GetUsers(usrStorage))
	r.POST("/api/user", api.Authenticate(), api.AddUser(usrStorage))
	r.PUT("/api/user", api.Authenticate(), api.UpdateUser(usrStorage))
	r.DELETE("/api/user/*username", api.Authenticate(), api.DeleteUser(usrStorage))

	r.Run()
}
