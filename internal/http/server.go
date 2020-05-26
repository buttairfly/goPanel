package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/http/router"
	"github.com/buttairfly/goPanel/internal/http/router/api"
)

// RunHTTPServer starts and runs a echo http server
func RunHTTPServer(logger *zap.Logger) {
	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	router.StaticRouter(e)
	api.Router(e)

	logger.Fatal("error in http server", zap.Error(e.Start(":8080")))
}
