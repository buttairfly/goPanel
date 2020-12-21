package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/weberror"
	"github.com/buttairfly/goPanel/pkg/version"
)

// GetVersions returns all program versions
func GetVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, version.GetVersions())
}

// GetVersionByName returns the version by name
func GetVersionByName(c echo.Context) error {
	name := c.Param("name")
	version, err := version.GetVersionByName(name)
	if err != nil {
		return weberror.NotFound("version", name)
	}
	return c.JSON(http.StatusOK, version)
}
