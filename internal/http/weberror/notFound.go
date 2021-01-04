package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NotFound returns an echo Error
func NotFound(entity string, value string, err error) error {
	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("NotFound: Entity '%s' with value %s error: %v", entity, value, err))
}
