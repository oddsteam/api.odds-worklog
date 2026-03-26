package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.odds.team/worklog/api.odds-worklog/api/file"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/login"
	sap_export_failure "gitlab.odds.team/worklog/api.odds-worklog/api/sap_export_failure"
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

// @title Odds-Worklog Example API
// @version 1.0
// @description This is a sample server odds-worklog server.
// @host http://worklog-dev.odds.team/api
// @BasePath /v1
func main() {
	_ = godotenv.Load()
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		log.Fatal("JWT_SIGNING_KEY environment variable is required")
	}

	session := mongo.Setup()
	defer session.Close()

	// Echo instance
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Middleware
	m := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte(jwtSigningKey),
	}

	r := e.Group("/v1")
	login.NewHttpHandler(r, session)
	r.Use(middleware.JWTWithConfig(m))

	// Handler
	user.NewHttpHandler(r, session)
	income.NewHttpHandler(r, session)
	file.NewHttpHandler(r, session)
	site.NewHttpHandler(r, session)
	sap_export_failure.NewHttpHandler(r, session)

	// Start server
	c := config.Config()
	e.Logger.Fatal(e.Start(c.APIPort))
}
