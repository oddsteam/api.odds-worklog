package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

func main() {
	session, err := mongo.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	// Echo instance
	e := echo.New()

	m := middleware.InitMiddleware()
	e.Use(m.CORS)

	// Routes
	e.GET("/", hello)

	user.NewHttpHandler(e, session)

	// Start server
	e.Logger.Fatal(e.Start(config.APIPort))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
