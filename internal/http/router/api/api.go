package api

import (
	"github.com/labstack/echo/v4"

	apiv1 "github.com/buttairfly/goPanel/internal/http/router/api/v1"
)

// Router changes the api router
func Router(e *echo.Echo) {
	g := e.Group("/api/v1")

	apiv1.Version(g)
	apiv1.Config(g)
}
