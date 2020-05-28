package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/config"
)

// GetMainConfig returns the MainConfig as JSON
func GetMainConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, config.GetMainConfig())
}

// GetConsumerConfig returns an array of TileConfigs as JSON
func GetConsumerConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, config.GetConsumerConfig())
}
