package apiv1

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Version adds all version routes
func Version(g *echo.Group) {
	versionGroup := g.Group("/version")
	versionGroup.GET("/all", controller.GetVersions).Name = "get-all-versions"
	versionGroup.GET("/name/:name", controller.GetVersionByName).Name = "get-version-by-name"
}
