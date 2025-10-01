package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port      string
    MongoURI  string
    MongoDB   string
    JWTSecret string
    JWTExpMin int
}

func New() *Config {
    exp := 60
    if v := os.Getenv("JWT_EXP_MIN"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            exp = n
        }
    }
    return &Config{
        Port:      getenv("PORT", "3001"),
        MongoURI:  getenv("MONGO_URI", "mongodb://localhost:27017"),
        MongoDB:   getenv("MONGO_DB", "authdb"),
        JWTSecret: getenv("JWT_SECRET", "supersecreto"),
        JWTExpMin: exp,
    }
}

func getenv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}
