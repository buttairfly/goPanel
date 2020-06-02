package apiv1

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Paint adds all paint routes
func Paint(g *echo.Group) {
	configGroup := g.Group("/paint")
	configGroup.GET("/pixel/FrameId/:frameId", controller.GetPixelColor).Name = "get-pixel-color"
	configGroup.PUT("/pixel/FrameId/:frameId", controller.SetPixelColor).Name = "set-pixel-color"

	configGroup.PUT("/fill/FrameId/:frameId", controller.SetFillColor).Name = "set-fill-color"

	// configGroup.GET("/frame/FrameId/:frameId", controller.GetFrame).Name = "get-frame-color"
	// configGroup.PUT("/frame/FrameId/:frameId", controller.SetFrame).Name = "put-frame-color"
}
