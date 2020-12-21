package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OutOfBounds returns an echo Error
func OutOfBounds(entity string, value string) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("OutOfBounds: Entity '%s' with value %s", entity, value))
}
