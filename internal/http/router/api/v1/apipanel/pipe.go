package apipanel

import (
	"github.com/buttairfly/goPanel/internal/http/controller"
	"github.com/labstack/echo/v4"
)

// Pipe adds all palettte routes
func Pipe(g *echo.Group) {
	configGroup := g.Group("/pipe")

	configGroup.GET("/all", controller.GetPanelPipes).Name = "get-all-panel-pipes"
	configGroup.GET("/id/:id", controller.GetPipeByID).Name = "get-pipe-by-id"

	configGroup.GET("/types/available", controller.GetAvailablePipeTypes).Name = "get-available-pipe-types"
	configGroup.GET("/types/reserved", controller.GetReservedPipeTypes).Name = "get-reserved-pipe-types"
}
