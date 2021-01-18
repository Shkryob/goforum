package handler

import (
	// "time"

	"github.com/shkryob/goforum/model"
	// "github.com/shkryob/goforum/utils"
	"github.com/labstack/echo/v4"
)

type postResponse struct {
	ID           uint    `json:"id"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
}

type singlePostResponse struct {
	Post *postResponse `json:"post"`
}

type postListResponse struct {
	Posts      []*postResponse `json:"postss"`
	PostsCount int                `json:"postsCount"`
}

func newPostResponse(c echo.Context, a *model.Post) *singlePostResponse {
	ar := new(postResponse)
	ar.ID = a.ID
	ar.Title = a.Title
	ar.Body = a.Body
	return &singlePostResponse{ar}
}

func newPostListResponse(posts []model.Post, count int) *postListResponse {
	r := new(postListResponse)
	r.Posts = make([]*postResponse, 0)
	for _, a := range posts {
		ar := new(postResponse)
		ar.Title = a.Title
		ar.Body = a.Body
		r.Posts = append(r.Posts, ar)
	}
	r.PostsCount = count
	return r
}