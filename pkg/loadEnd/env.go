package loadend

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warn("`.env` file not found, using existing environment variables")
			return
		}

		log.Fatalf("⚠️ Error loading .env file: %v", err)
	}
}
