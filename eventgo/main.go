package main

import (
    "log"

    "github.com/joho/godotenv"

    "github.com/tuusuario/eventgo/internal/di"
    "github.com/tuusuario/eventgo/internal/env"
    "github.com/tuusuario/eventgo/internal/rest"
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
    log.Println("Event server listening on", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}
