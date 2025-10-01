package main

import (
    "log"

    "github.com/joho/godotenv"

    "github.com/tuusuario/signupgo/internal/di"
    "github.com/tuusuario/signupgo/internal/env"
    "github.com/tuusuario/signupgo/internal/rest"
)

func main() {
    _ = godotenv.Load()

    cfg := env.New()
    inj, err := di.Build(cfg)
    if err != nil {
        log.Fatal("injector build error: ", err)
    }
    defer inj.Close()

    r := rest.NewRouter(inj)
    addr := ":" + cfg.Port
    log.Println("Signup server listening on", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}
