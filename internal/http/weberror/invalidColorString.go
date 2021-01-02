package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// InvalidColorString returns an echo Error
func InvalidColorString(entity string, id string, colorString string, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("InvalidColorString: %s: %s colorString %s error: %v", entity, id, colorString, err))
}
