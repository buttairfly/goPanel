package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/pkg/marshal"
)

// FramePoint is the
type FramePoint struct {
	FrameID string
	Point   marshal.Point
}

// GetPixelColor returns the PixelColor as color string
func GetPixelColor(c echo.Context) error {
	x, errX := strconv.Atoi(c.QueryParam("x"))
	if errX != nil {
		return errX
	}
	y, errY := strconv.Atoi(c.QueryParam("y"))
	if errY != nil {
		return errY
	}
	mp := marshal.Point{X: x, Y: y}
	// check if mp is within frame bounds

	fp := FramePoint{
		FrameID: c.Param("frameId"),
		Point:   mp,
	}
	zap.L().Sugar().Infof("%+v", fp)
	return c.JSON(http.StatusOK, fp)
}
