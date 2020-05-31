package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/marshal"
)

// FramePoint is the point in a specific frame
type FramePoint struct {
	FrameID string        `json:"frameId"`
	Point   marshal.Point `json:"point"`
}

// ColorAtFrame is a RGBA color at a FramePoint
type ColorAtFrame struct {
	FramePoint FramePoint `json:"framePoint"`
	Color      string     `json:"color"`
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
	currentFrame := device.GetLedDevice().GetCurrentFrame()
	if !fp.Point.ToImagePoint().In(currentFrame.Bounds()) {
		return fmt.Errorf("Point out of bounds of frame x %d y %d", x, y)
	}
	color := hardware.NewPixelFromColor(currentFrame.At(x, y))
	cf := ColorAtFrame{
		FramePoint: fp,
		Color:      color.Hex(),
	}
	return c.JSON(http.StatusOK, cf)
}

// SetPixelColor sets the PixelColor at frameId
func SetPixelColor(c echo.Context) error {
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
