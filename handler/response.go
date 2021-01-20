package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
)

type postResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	User  struct {
		Username string `json:"username"`
	} `json:"user"`
}

type singlePostResponse struct {
	Post *postResponse `json:"post"`
}

type postListResponse struct {
	Posts      []*postResponse `json:"posts"`
	PostsCount int             `json:"postsCount"`
}

func newPostResponse(c echo.Context, a *model.Post) *singlePostResponse {
	ar := new(postResponse)
	ar.ID = a.ID
	ar.Title = a.Title
	ar.Body = a.Body
	ar.User.Username = a.User.Username
	return &singlePostResponse{ar}
}

func newPostListResponse(posts []model.Post, count int) *postListResponse {
	r := new(postListResponse)
	r.Posts = make([]*postResponse, 0)
	for _, a := range posts {
		ar := new(postResponse)
		ar.ID = a.ID
		ar.Title = a.Title
		ar.Body = a.Body
		ar.User.Username = a.User.Username
		r.Posts = append(r.Posts, ar)
	}
	r.PostsCount = count
	return r
}

type commentResponse struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      struct {
		Username string `json:"username"`
	} `json:"user"`
}

type singleCommentResponse struct {
	Comment *commentResponse `json:"comment"`
}

type commentListResponse struct {
	Comments []commentResponse `json:"comments"`
}

func newCommentResponse(c echo.Context, cm *model.Comment) *singleCommentResponse {
	comment := new(commentResponse)
	comment.ID = cm.ID
	comment.Body = cm.Body
	comment.CreatedAt = cm.CreatedAt
	comment.UpdatedAt = cm.UpdatedAt
	comment.User.Username = cm.User.Username
	return &singleCommentResponse{comment}
}

func newCommentListResponse(c echo.Context, comments []model.Comment) *commentListResponse {
	r := new(commentListResponse)
	cr := commentResponse{}
	r.Comments = make([]commentResponse, 0)
	for _, i := range comments {
		cr.ID = i.ID
		cr.Body = i.Body
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.User.Username = i.User.Username

		r.Comments = append(r.Comments, cr)
	}
	return r
}

type userResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Token = utils.GenerateJWT(u.ID)
	return r
}
