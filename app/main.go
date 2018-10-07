package main

import (
	"net/http"

	"api.odds-worklog/app/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/userinfo", service.UserInfo)
	e.POST("/user", service.InsertUser)
	e.GET("/user", service.GetUser)
	e.DELETE("/user/:id", service.DeleteUser)
	e.PUT("/user/:id", service.UpdateUser)
	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
