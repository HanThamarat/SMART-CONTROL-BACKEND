package loadend

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load();

	if err != nil {
		log.Fatal("⚠️ Error loading .env file: %s", err);
	}
}