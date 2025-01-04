package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	// Load the .env file which is in the parent directory
	//err := godotenv.Load("../.env")
	err := godotenv.Load(".env")

	fmt.Println(err)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
