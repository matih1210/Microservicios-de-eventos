package env

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string

	AuthBaseURL   string
	HTTPTimeoutMs int
}

func New() *Config {
	return &Config{
		Port:          getenv("PORT", "3002"),
		MongoURI:      getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getenv("MONGO_DB", "eventsdb"),
		JWTSecret:     getenv("JWT_SECRET", "supersecreto"),
		AuthBaseURL:   getenv("AUTH_BASE_URL", "http://localhost:3001"),
		HTTPTimeoutMs: atoi(getenv("HTTP_TIMEOUT_MS", "2000")),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
