package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/weberror"
	"github.com/buttairfly/goPanel/internal/panel"
	"github.com/buttairfly/goPanel/pkg/palette"
)

// GetPalettes returns all panel palettes
func GetPalettes(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().GetMarshalledPalettes())
}

// GetPaletteByID returns the panel palette by id
func GetPaletteByID(c echo.Context) error {
	id := c.Param("id")
	palette, err := panel.GetPanel().GetMarshaledPaletteByID(id)
	if err != nil {
		return weberror.NotFound("palette", id)
	}
	return c.JSON(http.StatusOK, palette)
}

// PostColorAtPosToPaletteByID adds a new color to palette with id
func PostColorAtPosToPaletteByID(c echo.Context) error {
	id := c.Param("id")
	p, err := panel.GetPanel().GetPaletteByID(id)
	if err != nil {
		return weberror.NotFound("palette", id)
	}

	var paletteColor palette.ColorMarshal
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	errDecode := decoder.Decode(&paletteColor)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return weberror.BodyUnmarshal("palette.ColorMarshal", unmarshalErr)
		}
		return weberror.BodyNotFound(c, errDecode)
	}
	p.PutAt(paletteColor.Color, paletteColor.Pos)
	return c.JSON(http.StatusOK, p.ToMarshal())
}

// PutMoveColorAtPaletteByID moves a color fixpoint within the palette scale
func PutMoveColorAtPaletteByID(c echo.Context) error {
	id := c.Param("id")
	p, err := panel.GetPanel().GetPaletteByID(id)
	if err != nil {
		return weberror.NotFound("palette", id)
	}

	var paletteMove palette.ColorMoveMarshal
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	errDecode := decoder.Decode(&paletteMove)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return weberror.BodyUnmarshal("palette.ColorMarshal", unmarshalErr)
		}
		return weberror.BodyNotFound(c, errDecode)
	}

	errMove := p.MoveAt(paletteMove)
	if errMove != nil {
		return weberror.NotPossible("move not possible", errMove)
	}
	return c.JSON(http.StatusOK, p.ToMarshal())
}
