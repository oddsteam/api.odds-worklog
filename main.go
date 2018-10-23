package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

func main() {
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

	session, err := mongo.NewSession(&config)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	// Echo instance
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	m := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}

	// Handler
	user.NewHttpHandler(e, m, session)
	income.NewHttpHandler(e, m, session)

	// Start server
	e.Logger.Fatal(e.Start(config.APIPort))
}
