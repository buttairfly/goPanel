package weberror

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ColorConversion returns an echo Error
func ColorConversion(colorString string, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ColorConversion: '%s' has error %+v", colorString, err))
}
