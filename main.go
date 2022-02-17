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
			IsAdmin:  true,
		}
		err := usrStorage.Add(adminUser)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		err = usrStorage.SetPassword(adminUser.UserName, utility.EncondePassword(cfg.Admin.Password))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	sec := storageprovider.GetSecurityStorage(cfg.Provider)

	cryptoPrivateKey, cryptoPublicKey := utility.InitKeyPair(sec.LoadCryptingPair, sec.SaveCryptingPair)
	signinPrivateKey, signinPublicKey := utility.InitKeyPair(sec.LoadSigninPair, sec.SaveSigninPair)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/api/key/*path", api.Get(stg))
	router.GET("/api/keys/*pathSearch", api.GetKeys(stg))
	router.GET("/api/keys", api.GetAllKeys(stg))

	// Create or update key must be an authenticated call
	router.PUT("/api/key/*path", api.Authenticate(signinPublicKey, false), api.Put(stg, cryptoPublicKey))
	router.DELETE("/api/key/*path", api.Authenticate(signinPublicKey, false), api.Delete(stg))
	router.GET("api/key-enc/*path", api.Authenticate(signinPublicKey, false), api.GetEncrypted(stg, cryptoPrivateKey))

	router.POST("/api/token", api.Token(usrStorage, signinPrivateKey))

	router.GET("/api/user/:username", api.Authenticate(signinPublicKey, true), api.GetUser(usrStorage))
	router.GET("/api/users", api.Authenticate(signinPublicKey, true), api.GetUsers(usrStorage))
	router.POST("/api/user", api.Authenticate(signinPublicKey, true), api.AddUser(usrStorage))
	router.PUT("/api/user", api.Authenticate(signinPublicKey, true), api.UpdateUser(usrStorage))
	router.PATCH("/api/user/password/:username", api.Authenticate(signinPublicKey, true), api.SetPassword(usrStorage))
	router.DELETE("/api/user/*username", api.Authenticate(signinPublicKey, true), api.DeleteUser(usrStorage))

	router.Run()
}
