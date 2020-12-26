package weberror

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BodyUnmarshal returns an echo Error
func BodyUnmarshal(structString string, err *json.UnmarshalTypeError) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("BodyUnmarshal: could not unmarshal %s error: %v", structString, err))
}
