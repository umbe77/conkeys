package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"conkeys/api"
	"conkeys/storageprovider"
)


func main() {
    r := gin.Default()
	stg := storageprovider.GetKeyStorage("mongodb")

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
