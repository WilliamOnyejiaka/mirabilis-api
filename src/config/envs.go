package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Envs(key string) string {
	return map[string]string{
		"port":   os.Getenv("PORT"),
		"databaseUrl": os.Getenv("DATABASE_URL"),
		"databaseName": os.Getenv("DATABASE_NAME"),
		"environment": os.Getenv("ENVIRONMENT"),
		"secretKey": os.Getenv("SECRET_KEY"),
	}[key]
}
