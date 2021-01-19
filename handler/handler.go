package handler

import (
	"github.com/shkryob/goforum/model"
)

type Handler struct {
	userStore UserStore
	postStore PostStore
}

type PostStore interface {
	GetById(uint) (*model.Post, error)
	List(int, int) ([]model.Post, int, error)
	CreatePost(*model.Post) error
	UpdatePost(*model.Post) error
	DeletePost(*model.Post) error

	AddComment(*model.Post, *model.Comment) error
	GetCommentsByPostId(uint) ([]model.Comment, error)
	GetCommentByID(uint) (*model.Comment, error)
	DeleteComment(*model.Comment) error
}

type UserStore interface {
	GetByID(uint) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Create(*model.User) error
	Update(*model.User) error
}

func NewHandler(us UserStore, ps PostStore) *Handler {
	return &Handler{
		userStore: us,
		postStore: ps,
	}
}
