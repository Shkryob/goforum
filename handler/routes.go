package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	posts := v1.Group("/posts")

	// posts.GET("", h.Posts)
	posts.GET("/:id", h.GetPost)
}
