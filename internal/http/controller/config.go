package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/panel"
)

// GetMainConfig returns the MainConfig as JSON
func GetMainConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().GetMainConfig())
}

// GetConsumerConfig returns an array of TileConfigs as JSON
func GetConsumerConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().GetConsumerConfig())
}
