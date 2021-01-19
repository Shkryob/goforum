package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
)

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (request *userRegisterRequest) bind(context echo.Context, user *model.User) error {
	if err := context.Bind(request); err != nil {
		return err
	}
	if err := context.Validate(request); err != nil {
		return err
	}
	user.Username = request.User.Username
	user.Email = request.User.Email
	h, err := user.HashPassword(request.User.Password)
	if err != nil {
		return err
	}
	user.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type postCreateRequest struct {
	Post struct {
		Title string `json:"title" validate:"required"`
		Body  string `json:"body" validate:"required"`
	} `json:"post"`
}

func (request *postCreateRequest) bind(context echo.Context, a *model.Post) error {
	if err := context.Bind(request); err != nil {
		return err
	}
	if err := context.Validate(request); err != nil {
		return err
	}
	a.Title = request.Post.Title
	a.Body = request.Post.Body
	return nil
}

type postUpdateRequest struct {
	Post struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"post"`
}

func (request *postUpdateRequest) populate(post *model.Post) {
	request.Post.Title = post.Title
	request.Post.Body = post.Body
}

func (request *postUpdateRequest) bind(context echo.Context, a *model.Post) error {
	if err := context.Bind(request); err != nil {
		return err
	}
	if err := context.Validate(request); err != nil {
		return err
	}
	a.Title = request.Post.Title
	a.Body = request.Post.Body
	return nil
}

type createCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}

func (request *createCommentRequest) bind(context echo.Context, cm *model.Comment) error {
	if err := context.Bind(request); err != nil {
		return err
	}
	if err := context.Validate(request); err != nil {
		return err
	}
	cm.Body = request.Comment.Body
	return nil
}
