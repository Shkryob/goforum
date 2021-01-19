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

// GetPost godoc
// @Summary Get a post
// @Description Get a post. Auth not required
// @ID get-post
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the post to get"
// @Success 200 {object} singlePostResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /posts/{id} [get]
func (handler *Handler) GetPost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)
	post, err := handler.postStore.GetById(id)

	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}

	return utils.ResponseByContentType(context, http.StatusOK, newPostResponse(context, post))
}

// GetPosts godoc
// @Summary Get recent posts globally
// @Description Get most recent posts globally. Use query parameters to filter results. Auth is optional
// @ID get-posts
// @Tags post
// @Accept  json
// @Produce  json
// @Param limit query integer false "Limit number of posts returned (default is 20)"
// @Param offset query integer false "Offset/skip number of posts (default is 0)"
// @Success 200 {object} postListResponse
// @Failure 500 {object} utils.Error
// @Router /posts [get]
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

	return utils.ResponseByContentType(context, http.StatusOK, newPostListResponse(posts, count))
}

// CreatePost godoc
// @Summary Create a post
// @Description Create a post. Auth is require
// @ID create-post
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body postCreateRequest true "Post to create"
// @Success 201 {object} singlePostResponse
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /posts [post]
func (handler *Handler) CreatePost(context echo.Context) error {
	var post model.Post

	req := &postCreateRequest{}
	if err := req.bind(context, &post); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}

	post.UserID = userIDFromToken(context)

	err := handler.postStore.CreatePost(&post)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusCreated, newPostResponse(context, &post))
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post. Auth is required
// @ID update-post
// @Tags post
// @Accept  json
// @Produce  json
// @Param ID path string true "ID of the post to update"
// @Param post body postUpdateRequest true "Post to update"
// @Success 200 {object} singlePostResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /posts/{id} [put]
func (handler *Handler) UpdatePost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}

	req := &postUpdateRequest{}
	req.populate(post)

	if err := req.bind(context, post); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = handler.postStore.UpdatePost(post); err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusOK, newPostResponse(context, post))
}

// DeletPost godoc
// @Summary Delete a post
// @Description Delete a post. Auth is required
// @ID delete-post
// @Tags post
// @Accept  json
// @Produce  json
// @Param ID path string true "ID of the post to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
func (handler *Handler) DeletePost(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}

	err = handler.postStore.DeletePost(post)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusOK, map[string]interface{}{"result": "ok"})
}

// AddComment godoc
// @Summary Create a comment for a post
// @Description Create a comment for a post. Auth is required
// @ID add-comment
// @Tags comment
// @Accept  json
// @Produce  json
// @Param ID path string true "ID of the post that you want to create a comment for"
// @Param comment body createCommentRequest true "Comment you want to create"
// @Success 201 {object} singleCommentResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /posts/{id}/comments [post]
func (handler *Handler) AddComment(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	post, err := handler.postStore.GetById(id)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	if post == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}

	var cm model.Comment

	req := &createCommentRequest{}
	if err := req.bind(context, &cm); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}

	cm.UserID = userIDFromToken(context)

	if err = handler.postStore.AddComment(post, &cm); err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusCreated, newCommentResponse(context, &cm))
}

// GetComments godoc
// @Summary Get the comments for a post
// @Description Get the comments for a post. Auth is optional
// @ID get-comments
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the post that you want to get comments for"
// @Success 200 {object} commentListResponse
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /posts/{id}/comments [get]
func (handler *Handler) GetComments(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("post_id"), 10, 32)
	id := uint(id64)

	cm, err := handler.postStore.GetCommentsByPostId(id)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusOK, newCommentListResponse(context, cm))
}

// DeleteComment godoc
// @Summary Delete a comment for a post
// @Description Delete a comment for a post. Auth is required
// @ID delete-comments
// @Tags comment
// @Accept  json
// @Produce  json
// @Param ID path string true "ID of the post that you want to delete a comment for"
// @Param id path integer true "ID of the comment you want to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /posts/{id}/comments/{id} [delete]
func (handler *Handler) DeleteComment(context echo.Context) error {
	id64, err := strconv.ParseUint(context.Param("comment_id"), 10, 32)
	id := uint(id64)

	if err != nil {
		return utils.ResponseByContentType(context, http.StatusBadRequest, utils.NewError(err))
	}

	cm, err := handler.postStore.GetCommentByID(id)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	if cm == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}

	if cm.UserID != userIDFromToken(context) {
		return utils.ResponseByContentType(context, http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}

	if err := handler.postStore.DeleteComment(cm); err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}

	return utils.ResponseByContentType(context, http.StatusOK, map[string]interface{}{"result": "ok"})
}
