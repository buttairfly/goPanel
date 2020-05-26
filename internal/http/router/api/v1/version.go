package apiv1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/pkg/version"
)

// GetVersions returns
func getVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, version.Versions)
}

// Version adds all version routes
func Version(g *echo.Echo) {
	g.GET("/versions", getVersions).name = "get-versions"
}
