package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"conkeys/api"
	"conkeys/config"
	"conkeys/frontend"
	"conkeys/storage"
	"conkeys/storage/postgres"
	"conkeys/utility"
)

func main() {
    cfg := config.GetConfig()

    engine := utility.NewEngine("./frontend/views")

    app := fiber.New(fiber.Config{
        Views: engine,
    })

    if cfg.Postgres.ConnectionUri == "" {
        log.Fatal("No POstgres connection uri defined")
    }
    connectionUri := cfg.Postgres.ConnectionUri

    db, err := sql.Open("postgres", connectionUri)
    if err != nil {
        log.Fatal(err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    stg := postgres.NewKeyStorage(db)
    usrStorage := postgres.NewUserStorage(db)

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

    sec := postgres.NewSecStorage(db)

    cryptoPrivateKey, cryptoPublicKey := utility.InitKeyPair(sec.LoadCryptingPair, sec.SaveCryptingPair)
    signinPrivateKey, signinPublicKey := utility.InitKeyPair(sec.LoadSigninPair, sec.SaveSigninPair)

    frontend.InitControllers(app)

    app.Get("/ping", func(c *fiber.Ctx) error {
        c.Status(fiber.StatusOK)
        return c.JSON(fiber.Map{
            "message": "pong",
        })
    })

    app.Get("/api/key/*", api.Get(stg))
    app.Get("/api/keys/*", api.GetKeys(stg))
    app.Get("/api/keys", api.GetAllKeys(stg))

    // Create or update key must be an authenticated call
    app.Put("/api/key/*", api.Authenticate(signinPublicKey, false), api.Put(stg, cryptoPublicKey))
    app.Delete("/api/key/*", api.Authenticate(signinPublicKey, false), api.Delete(stg))
    app.Get("api/key-enc/*", api.Authenticate(signinPublicKey, false), api.GetEncrypted(stg, cryptoPrivateKey))

    app.Post("/api/token", api.Token(usrStorage, signinPrivateKey))

    app.Get("/api/user/:username", api.Authenticate(signinPublicKey, true), api.GetUser(usrStorage))
    app.Get("/api/users", api.Authenticate(signinPublicKey, true), api.GetUsers(usrStorage))
    app.Post("/api/user", api.Authenticate(signinPublicKey, true), api.AddUser(usrStorage))
    app.Put("/api/user", api.Authenticate(signinPublicKey, true), api.UpdateUser(usrStorage))
    app.Patch("/api/user/password/:username", api.Authenticate(signinPublicKey, true), api.SetPassword(usrStorage))
    app.Delete("/api/user/:username", api.Authenticate(signinPublicKey, true), api.DeleteUser(usrStorage))

    app.Listen(":8080")
}
