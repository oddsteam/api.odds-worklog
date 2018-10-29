package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	_ "gitlab.odds.team/worklog/api.odds-worklog/docs"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/login"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

func main() {
	c := config.Config()
	// Setup Mongo
	session, err := mongo.NewSession(c)
	if err != nil {
		log.Fatal(err.Error())
	}
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
		SigningKey: []byte("GmkZGF3CmpZNs88dLvbV"),
	}

	r := e.Group("/v1")
	r.GET("/swagger/*", echoSwagger.WrapHandler)
	login.NewHttpHandler(r, session)
	r.Use(middleware.JWTWithConfig(m))

	// Handler
	user.NewHttpHandler(r, session)
	income.NewHttpHandler(r, session)

	// Start server
	e.Logger.Fatal(e.Start(c.APIPort))
}
