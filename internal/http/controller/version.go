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
