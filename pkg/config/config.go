package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string
	AppPort    string
	JWTSecret  string
}

// LoadConfig reads .env file and returns a Config struct
func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, falling back to system env vars")
	}

	cfg := Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		AppPort:    getEnv("APP_PORT", "3000"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Shanghai"),
		JWTSecret:  getEnv("JWT_SECRET", "secret-key"),
	}

	return cfg, nil
}

// DSN builds the PostgreSQL connection string from config
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode, c.DBTimeZone,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
