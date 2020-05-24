package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// RunHTTPServer starts and runs a echo http server
func RunHTTPServer(logger *zap.Logger) {
	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.File("/favicon.ico", "images/favicon.ico")
	logger.Fatal("error in http server", zap.Error(e.Start(":8080")))
}
