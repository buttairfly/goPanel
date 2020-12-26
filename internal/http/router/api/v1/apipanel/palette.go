package apipanel

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Palette adds all palettte routes
func Palette(g *echo.Group) {
	configGroup := g.Group("/palette")

	configGroup.GET("/all", controller.GetPalettes).Name = "get-all-palettes"
	configGroup.GET("/id/:id", controller.GetPaletteByID).Name = "get-palette-by-id"
	configGroup.POST("/id/:id/color/set", controller.PostColorAtPosToPaletteByID).Name = "set-color-at-pos-to-palette-by-id"
	configGroup.PUT("/id/:id/color/move", controller.PutMoveColorAtPaletteByID).Name = "move-color-at-palette-by-id"

	// configGroup.POST("/id/:id", controller.PostPaletteById).Name = "new-palette-by-id"
}
