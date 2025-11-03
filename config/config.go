package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Getenv(key string) string {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Printf("failed to load env file %v", err)
	}

	return os.Getenv(key)
}
