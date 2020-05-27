package apiv1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/config"
)

// GetVersions returns
func getMainConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, config.GetMainConfig())
}

// Config adds all config routes
func Config(g *echo.Group) {
	configGroup := g.Group("/config")
	configGroup.GET("/all", getMainConfig).Name = "get-main-config"
}
