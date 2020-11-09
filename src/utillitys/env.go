package utillitys

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//  Getter environment from .env.
func GetENV(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}