package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error in loading .env file")
	}
}
