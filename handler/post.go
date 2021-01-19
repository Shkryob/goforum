package handler

import (
	// "errors"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
)

func (handler *Handler) GetPost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)
	post, err := handler.postStore.GetById(id)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return context.JSON(http.StatusNotFound, utils.NotFound())
	}

	return context.JSON(http.StatusOK, newPostResponse(context, post))
}

func (handler *Handler) GetPosts(context echo.Context) error {
	var (
		posts []model.Post
		count int
	)

	offset, err := strconv.Atoi(context.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(context.QueryParam("limit"))
	if err != nil {
		limit = 20
	}

	posts, count, err = handler.postStore.List(offset, limit)

	return context.JSON(http.StatusOK, newPostListResponse(posts, count))
}

func (handler *Handler) CreatePost(context echo.Context) error {
	var post model.Post

	req := &postCreateRequest{}
	if err := req.bind(context, &post); err != nil {
		return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	post.UserID = userIDFromToken(context)

	err := handler.postStore.CreatePost(&post)
	if err != nil {
		return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return context.JSON(http.StatusCreated, newPostResponse(context, &post))
}

func (handler *Handler) UpdatePost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return context.JSON(http.StatusNotFound, utils.NotFound())
	}

	req := &postUpdateRequest{}
	req.populate(post)

	if err := req.bind(context, post); err != nil {
		return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = handler.postStore.UpdatePost(post); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return context.JSON(http.StatusOK, newPostResponse(context, post))
}

func (handler *Handler) DeletePost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return context.JSON(http.StatusNotFound, utils.NotFound())
	}

	err = handler.postStore.DeletePost(post)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return context.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (handler *Handler) AddComment(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return context.JSON(http.StatusNotFound, utils.NotFound())
	}

	var cm model.Comment

	req := &createCommentRequest{}
	if err := req.bind(context, &cm); err != nil {
		return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	cm.UserID = userIDFromToken(context)

	if err = handler.postStore.AddComment(post, &cm); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return context.JSON(http.StatusCreated, newCommentResponse(context, &cm))
}

func (handler *Handler) GetComments(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	cm, err := handler.postStore.GetCommentsByPostId(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return context.JSON(http.StatusOK, newCommentListResponse(context, cm))
}

func (handler *Handler) DeleteComment(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("comment_id"), 10, 32)
	id := uint(id64)

	if err != nil {
		return context.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	cm, err := handler.postStore.GetCommentByID(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if cm == nil {
		return context.JSON(http.StatusNotFound, utils.NotFound())
	}

	if cm.UserID != userIDFromToken(context) {
		return context.JSON(http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}

	if err := handler.postStore.DeleteComment(cm); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return context.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}
