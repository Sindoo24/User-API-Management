package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	ServerPort       string
	LogLevel         string
	JWTSecret        string
	JWTExpiry        time.Duration
	CookieSecure     bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse JWT expiry in hours, default to 24 hours
	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		expiryHours = 24
	}

	// Parse cookie secure flag, default to true (secure in production)
	cookieSecure := getEnv("COOKIE_SECURE", "true") == "true"

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:    time.Duration(expiryHours) * time.Hour,
		CookieSecure: cookieSecure,
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
