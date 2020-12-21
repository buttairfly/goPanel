package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Conversion returns an echo Error
func Conversion(entity string, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Conversion: '%s' has error %+v", entity, err))
}
