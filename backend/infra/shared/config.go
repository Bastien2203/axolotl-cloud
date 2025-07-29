package shared

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var expectedEnvVars = []string{
	"HTTP_PORT",
	"VOLUMES_PATH_HOST",
	"VOLUMES_PATH_CONTAINER",
	"DATABASE_PATH",
}

func LoadEnv() error {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	for _, envVar := range expectedEnvVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("environment variable %s is not set", envVar)
		}
	}
	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
