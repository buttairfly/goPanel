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
	FrameID string        `json:"frameId"`
	Point   marshal.Point `json:"point,omitempty"`
	Color   string        `json:"color"`
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
	frame := device.GetLedDevice().GetCurrentFrame()
	if !mp.ToImagePoint().In(frame.Bounds()) {
		return fmt.Errorf("Point out of bounds of frame x %d y %d", x, y)
	}
	color := hardware.NewPixelFromColor(frame.At(x, y))
	cf := ColorAtFrame{
		FrameID: frameID,
		Point:   mp,
		Color:   color.Hex(),
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
	frame := device.GetLedDevice().GetCurrentFrame()
	if !mp.ToImagePoint().In(frame.Bounds()) {
		return fmt.Errorf("Point out of bounds of frame x %d y %d", x, y)
	}

	color, errColor := hardware.NewPixelFromHex(c.QueryParam("color"))
	if errColor != nil {
		return errColor
	}

	// set pixel
	frame.Set(x, y, color)

	cf := ColorAtFrame{
		FrameID: frameID,
		Point:   mp,
		Color:   color.Hex(),
	}

	return c.JSON(http.StatusOK, cf)
}

// SetFillColor sets the PixelColor at frameId
func SetFillColor(c echo.Context) error {
	frameID := c.Param("frameId")
	// TODO: use frameID instead of currentFrame
	frame := device.GetLedDevice().GetCurrentFrame()

	color, errColor := hardware.NewPixelFromHex(c.QueryParam("color"))
	if errColor != nil {
		return errColor
	}

	// set pixel
	frame.Fill(color)

	cf := ColorAtFrame{
		FrameID: frameID,
		Color:   color.Hex(),
	}

	return c.JSON(http.StatusOK, cf)
}
