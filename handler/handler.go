package handler

import (
	"github.com/shkryob/goforum/model"
)

type Handler struct {
	userStore    UserStore
	postStore PostStore
}

type PostStore interface {
	GetById(uint) (*model.Post, error)
}

type UserStore interface {
}

func NewHandler(us UserStore, ps PostStore) *Handler {
	return &Handler{
		userStore:    us,
		postStore: ps,
	}
}