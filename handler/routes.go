package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/router/middleware"
	"github.com/shkryob/goforum/utils"
)

func (handler *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	guestUsers := v1.Group("/users")
	guestUsers.POST("", handler.SignUp)
	guestUsers.POST("/login", handler.Login)

	user := v1.Group("/user", jwtMiddleware)
	user.GET("", handler.CurrentUser)

	posts := v1.Group("/posts")

	posts.GET("", handler.GetPosts)
	posts.GET("/:post_id", handler.GetPost)
	posts.POST("", handler.CreatePost, jwtMiddleware)
	posts.PUT("/:post_id", handler.UpdatePost, jwtMiddleware)
	posts.DELETE("/:post_id", handler.DeletePost, jwtMiddleware)

	comments := v1.Group("/posts/:post_id/comments")

	comments.GET("", handler.GetComments)
	comments.POST("", handler.AddComment, jwtMiddleware)
	comments.DELETE("/:comment_id", handler.DeleteComment, jwtMiddleware)
}
