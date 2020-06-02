package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/pkg/version"
)

// GetVersions returns all program versions
func GetVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, version.GetVersions())
}

// GetVersionByName returns the version by name
func GetVersionByName(c echo.Context) error {
	version, err := version.GetVersionByName(c.Param("name"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, version)
}
