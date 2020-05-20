package http

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// RunHTTPServer starts and runs a echo http server
func RunHTTPServer(wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()

	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.File("/favicon.ico", "images/favicon.png")
	logger.Fatal("error in http server", zap.Error(e.Start(":8080")))
}
