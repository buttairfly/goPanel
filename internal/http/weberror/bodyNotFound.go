package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BodyNotFound returns an echo Error
func BodyNotFound(c echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("BodyNotFound: Route '%s' with method %s error: %v", c.Request().RequestURI, c.Request().Method, err))
}
