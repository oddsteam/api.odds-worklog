package main

import (
	"log"
	"net/http"
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
	// Routes
	e.GET("/", hello)

	user.NewHttpHandler(e, session)
	income.NewHttpHandler(e, session)
	// Start server
	e.Logger.Fatal(e.Start(config.APIPort))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
