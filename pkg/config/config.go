package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func Config() *models.Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cp, _ := strconv.Atoi(os.Getenv("MONGO_DB_CONECTION_POOL"))
	config := models.Config{
		os.Getenv("MONGO_DB_HOST"),
		os.Getenv("MONGO_DB_NAME"),
		cp,
		os.Getenv("API_PORT"),
	}
	return &config
}
