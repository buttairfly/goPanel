package apiv1

import (
	"github.com/buttairfly/goPanel/internal/http/router/api/v1/apipanel"
	"github.com/labstack/echo/v4"
)

// Panel adds all panel routes
func Panel(g *echo.Group) {
	panelGroup := g.Group("/panel")

	apipanel.Pipe(panelGroup)
	apipanel.Palette(panelGroup)
}
