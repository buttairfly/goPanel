package apiv1

import (
	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/controller"
)

// Config adds all config routes
func Config(g *echo.Group) {
	configGroup := g.Group("/config")
	configGroup.GET("/main", controller.GetMainConfig).Name = "get-main-config"
	configGroup.GET("/consumer", controller.GetConsumerConfig).Name = "get-consumer-config"
}
