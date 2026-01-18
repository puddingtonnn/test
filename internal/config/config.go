package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN   string
	AppPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbUser := getEnv("POSTGRES_USER", "postgres")
	dbPass := getEnv("POSTGRES_PASSWORD", "secret")
	dbName := getEnv("POSTGRES_DB", "testdb")
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")

	appPort := getEnv("HTTP_PORT", "9090")

	dbdsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	return &Config{
		DBDSN:   dbdsn,
		AppPort: appPort,
	}
}

// getEnv возвращает значение переменной окружения или дефолтное значение
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
