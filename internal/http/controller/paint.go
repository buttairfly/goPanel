package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/marshal"
)

// ColorAtFrame is a color at a FramePoint
type ColorAtFrame struct {
	FrameID    string     `json:"frameId"`
	PointColor PointColor `json:"pointColor"`
}

// PointColor is a color at a point
type PointColor struct {
	Point marshal.Point `json:"point"`
	Color string        `json:"color"`
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

	frameID := c.Param("frameId")
	// TODO: use frameID instead of currentFrame
	currentFrame := device.GetLedDevice().GetCurrentFrame()
	if !mp.ToImagePoint().In(currentFrame.Bounds()) {
		return fmt.Errorf("Point out of bounds of frame x %d y %d", x, y)
	}
	color := hardware.NewPixelFromColor(currentFrame.At(x, y))
	cf := ColorAtFrame{
		FrameID: frameID,
		PointColor: PointColor{
			Point: mp,
			Color: color.Hex(),
		},
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

	frameID := c.Param("frameId")
	// TODO: use frameID instead of currentFrame
	currentFrame := device.GetLedDevice().GetCurrentFrame()
	if !mp.ToImagePoint().In(currentFrame.Bounds()) {
		return fmt.Errorf("Point out of bounds of frame x %d y %d", x, y)
	}

	color, errColor := hardware.NewPixelFromHex(c.QueryParam("color"))
	if errColor != nil {
		return errColor
	}

	// set pixel
	currentFrame.Set(x, y, color)

	cf := ColorAtFrame{
		FrameID: frameID,
		PointColor: PointColor{
			Point: mp,
			Color: color.Hex(),
		},
	}

	return c.JSON(http.StatusOK, cf)
}
