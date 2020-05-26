package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StaticRouter surfs all static content
func StaticRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.File("/favicon.ico", "images/favicon.ico")
}
