package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NotPossible returns an echo Error
func NotPossible(description string, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("NotPossible: %s: %v", description, err))
}
