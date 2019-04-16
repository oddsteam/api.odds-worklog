package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func Config() *models.Config {
	godotenv.Load()

	cp, _ := strconv.Atoi(os.Getenv("MONGO_DB_CONECTION_POOL"))
	config := models.Config{
		os.Getenv("MONGO_DB_HOST"),
		os.Getenv("MONGO_DB_NAME"),
		cp,
		os.Getenv("API_PORT"),
		os.Getenv("MONGO_DB_USERNAME"),
		os.Getenv("MONGO_DB_PASSWORD"),
	}
	return &config
}
