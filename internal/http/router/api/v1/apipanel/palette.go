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
	// configGroup.POST("/id/:id/set/color/:color/pos/:pos", controller.PostColorAtPosToPaletteById).Name = "post-set-color-at-pos-to-palette-by-id"
	// configGroup.PUT("/id/:id/move/from-pos/:from-pos/to-pos/:to-pos", controller.PutMoveColorAtPaletteById).Name = "put-move-color-at-palette-by-id"

	// configGroup.POST("/id/:id", controller.PostPaletteById).Name = "post-palette-by-id"
}
