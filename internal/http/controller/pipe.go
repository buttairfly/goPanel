package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/buttairfly/goPanel/internal/http/weberror"
	"github.com/buttairfly/goPanel/internal/panel"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

// GetPipes returns all panel pipes
func GetPipes(c echo.Context) error {
	return c.JSON(http.StatusOK, panel.GetPanel().GetFramePipeline().MarshalFramePipeline())
}

// GetPipeByID returns the panel pipe by id
func GetPipeByID(c echo.Context) error {
	id := c.Param("id")
	pipe, err := panel.GetPanel().GetFramePipeline().GetPipeByID(pipepart.ID(id))
	if err != nil {
		return weberror.NotFound("pipe", id)
	}
	return c.JSON(http.StatusOK, pipe.Marshal())
}
