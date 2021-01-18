package handler

import (
	// "errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	// "github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
)

func (h *Handler) GetPost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	a, err := h.postStore.GetById(uint(id))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newPostResponse(c, a))
}
