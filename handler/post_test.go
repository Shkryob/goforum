package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/router"
	"github.com/shkryob/goforum/router/middleware"
	"github.com/shkryob/goforum/utils"
	"github.com/stretchr/testify/assert"
)

func TestListPostsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	e := router.New()
	req := httptest.NewRequest(echo.GET, "/api/posts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetPosts(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var aa postListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aa)
		assert.NoError(t, err)
		assert.Equal(t, 2, aa.PostsCount)
	}
}

func TestGetPostsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	req := httptest.NewRequest(echo.GET, "/api/posts/:post_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	assert.NoError(t, h.GetPost(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singlePostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), a.Post.ID)
	}
}

func TestCreatePostsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"post":{"title":"post2", "body":"post2"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/api/posts", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.CreatePost(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var a singlePostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "post2", a.Post.Title)
		assert.Equal(t, "user1", a.Post.User.Username)
	}
}

func TestUpdatePostsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"post":{"title":"post1 part 2"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.PUT, "/api/posts/:post_id", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.UpdatePost(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singlePostResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "post1 part 2", a.Post.Title)
	}
}

func TestDeletePostCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/posts/:post_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.DeletePost(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetCommentsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/posts/:post_id/comments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(2)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id/comments")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetComments(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cc commentListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cc)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(cc.Comments))
	}
}

func TestAddCommentCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"comment":{"body":"post1 comment2 by user2"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/api/posts/:post_id/comments", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(2)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id/comments")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.AddComment(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var c singleCommentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, "post1 comment2 by user2", c.Comment.Body)
		assert.Equal(t, "user2", c.Comment.User.Username)
	}
}

func TestDeleteCommentCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/posts/:post_id/comments/:comment_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/posts/:post_id/comments/:comment_id")
	c.SetParamNames("post_id")
	c.SetParamValues("1")
	c.SetParamNames("comment_id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.DeleteComment(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
