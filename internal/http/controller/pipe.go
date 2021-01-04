package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/weberror"
	"github.com/buttairfly/goPanel/internal/panel"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

// GetPanelPipes returns all panel pipes
func GetPanelPipes(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().Marshal())
}

// GetPipeByID returns the panel pipe by id
func GetPipeByID(c echo.Context) error {
	id := c.Param("id")
	pipe, err := panel.GetPanel().GetPipeByID(pipepart.ID(id))
	if err != nil {
		return weberror.NotFound("pipe", id, err)
	}
	return c.JSON(http.StatusOK, pipe.Marshal())
}

// GetReservedPipeTypes returns all reserved pipe types
func GetReservedPipeTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, pipepart.GetReservedPipeTypes())
}

// GetAvailablePipeTypes returns all available pipe types
func GetAvailablePipeTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, pipepart.GetPipeTypes())
}
