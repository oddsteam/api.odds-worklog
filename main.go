package main

import (
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/api/file"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/login"
	"gitlab.odds.team/worklog/api.odds-worklog/api/reminder"
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/worker"
)

// @title Odds-Worklog Example API
// @version 1.0
// @description This is a sample server odds-worklog server.
// @host http://worklog-dev.odds.team/api
// @BasePath /v1
func main() {
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
		SigningKey: []byte("GmkZGF3CmpZNs88dLvbV"),
	}

	r := e.Group("/v1")
	login.NewHttpHandler(r, session)
	r.Use(middleware.JWTWithConfig(m))

	// Handler
	user.NewHttpHandler(r, session)
	income.NewHttpHandler(r, session)
	reminder.NewHttpHandler(r, session)
	file.NewHttpHandler(r, session)
	site.NewHttpHandler(r, session)

	r = e.Group("/v2")
	r.Use(middleware.JWTWithConfig(m))
	income.NewHttpHandler2(r, session)

	reminderRepo := reminder.NewRepository(session)
	s, err := reminderRepo.GetReminder()
	if err != nil {
		log.Println(err)
	} else {
		worker.StartWorker(s)
	}
	// Start server
	c := config.Config()
	e.Logger.Fatal(e.Start(c.APIPort))
}
