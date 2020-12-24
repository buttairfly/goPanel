package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/weberror"
	"github.com/buttairfly/goPanel/internal/panel"
)

// GetPalettes returns all panel palettes
func GetPalettes(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().GetPalettes())
}

// GetPaletteByID returns the panel palette by id
func GetPaletteByID(c echo.Context) error {
	id := c.Param("id")
	palette, err := panel.GetPanel().GetPaletteByID(id)
	if err != nil {
		return weberror.NotFound("palette", id)
	}
	return c.JSON(http.StatusOK, palette)
}
