package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"conkeys/api"
	"conkeys/config"
	"conkeys/storageprovider"
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
	r.PUT("/api/key/*path", api.Put(stg))

    r.Run()
}
