package main

import (
	"conkeys/api"
	"conkeys/config"
	"conkeys/storageprovider"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.Provider)
	r := gin.Default()
	stg := storageprovider.GetKeyStorage(cfg.Provider)

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

	r.Run()
}
