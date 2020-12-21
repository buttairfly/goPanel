package http

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/http/router"
	"github.com/buttairfly/goPanel/internal/http/router/api"
)

// RunHTTPServer starts and runs a echo http server
func RunHTTPServer(cancelCtx context.Context, wg *sync.WaitGroup, gracePeriod time.Duration, logger *zap.Logger) {
	defer wg.Done()
	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	router.StaticRouter(e)
	api.Router(e)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err := e.Start(":8080"); err != nil {
			logger.Info("shutting down the server")
		}
	}(wg)

	// wait until interrupt comes
	<-cancelCtx.Done()
	timeoutContext, cancel := context.WithTimeout(cancelCtx, gracePeriod)
	defer cancel()
	if err := e.Shutdown(timeoutContext); err != nil {
		e.Logger.Error(err)
	}
	logger.Info("Server ended")
}
