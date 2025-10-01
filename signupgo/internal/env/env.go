package env

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	MongoURI    string
	MongoDB     string
	JWTSecret   string
	AuthBaseURL string

	EventBaseURL  string
	HTTPTimeoutMs int
}

func New() *Config {
	return &Config{
		Port:          getenv("PORT", "3003"),
		MongoURI:      getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getenv("MONGO_DB", "signupdb"),
		JWTSecret:     getenv("JWT_SECRET", "supersecreto"),
		EventBaseURL:  getenv("EVENT_BASE_URL", "http://localhost:3002"), // Event service base URL, se usa para comunicarse con el servicio de eventos
		AuthBaseURL:   getenv("AUTH_BASE_URL", "http://localhost:3001"),  // Auth service base URL, se usa para comunicarse con el servicio de autenticación
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
