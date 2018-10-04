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
	e.POST("/insertUser", service.InsertUser)
	e.GET("/getUser", service.GetUser)
	e.DELETE("/delete/:id", service.DeleteUser)
	e.PUT("/update/:id", service.UpdateUser)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
