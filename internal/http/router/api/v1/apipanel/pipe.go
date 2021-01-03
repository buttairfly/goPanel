package apipanel

import (
	"github.com/buttairfly/goPanel/internal/http/controller"
	"github.com/labstack/echo/v4"
)

// Pipe adds all palettte routes
func Pipe(g *echo.Group) {
	configGroup := g.Group("/pipe")

	configGroup.GET("/all", controller.GetPipes).Name = "get-all-pipes"
	configGroup.GET("/id/:id", controller.GetPipeByID).Name = "get-pipe-by-id"
}
