package apiv1

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Paint adds all paint routes
func Paint(g *echo.Group) {
	configGroup := g.Group("/paint")
	configGroup.GET("/pixel/frameId/:frameId", controller.GetPixelColor).Name = "get-pixel-color"
	configGroup.PUT("/pixel/frameId/:frameId", controller.SetPixelColor).Name = "set-pixel-color"

	configGroup.PUT("/fill/frameId/:frameId", controller.SetFillColor).Name = "set-fill-color"

	// configGroup.GET("/frame/frameId/:frameId", controller.GetFrame).Name = "get-frame-color"
	// configGroup.PUT("/frame/frameId/:frameId", controller.SetFrame).Name = "put-frame-color"
}
