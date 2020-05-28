package apiv1

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Paint adds all paint routes
func Paint(g *echo.Group) {
	configGroup := g.Group("/paint")
	configGroup.GET("/pixel/FrameId/:frameId", controller.GetPixelColor).Name = "get-pixel-color"
	// configGroup.PUT("/pixel/FrameId/:frameId", controller.PutPixelColor()).Name = "put-pixel-color"

	// configGroup.PUT("/fill/FrameId/:frameId", controller.PutFillColor()).Name = "put-fill-color"

	// configGroup.GET("/frame/FrameId/:frameId", controller.GetFrame()).Name = "get-frame-color"
	// configGroup.PUT("/frame/FrameId/:frameId", controller.PutFrame()).Name = "put-frame-color"
}
