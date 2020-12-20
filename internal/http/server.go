package http

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/http/router"
	"github.com/buttairfly/goPanel/internal/http/router/api"
)

// RunHTTPServer starts and runs a echo http server
func RunHTTPServer(cancelCtx context.Context, gracePeriod time.Duration, logger *zap.Logger) {
	e := echo.New()
	// Root level middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.HideBanner = true
	router.StaticRouter(e)
	api.Router(e)
	go func() {
		if err := e.Start(":8080"); err != nil {
			logger.Info("shutting down the server")
		}
	}()

	timeoutContext, cancel := context.WithTimeout(cancelCtx, gracePeriod)
	defer cancel()
	if err := e.Shutdown(timeoutContext); err != nil {
		e.Logger.Fatal(err)
	}
}
