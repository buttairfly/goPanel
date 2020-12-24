package api

import (
	"github.com/labstack/echo/v4"

	apiv1 "github.com/buttairfly/goPanel/internal/http/router/api/v1"
)

// Router changes the api router
func Router(e *echo.Echo) {
	g := e.Group("/api/v1")

	apiv1.Config(g)
	apiv1.Panel(g)
	apiv1.Panel(g)
	apiv1.Version(g)
}
