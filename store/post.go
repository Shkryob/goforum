package store

import (
	"github.com/jinzhu/gorm"
	"github.com/shkryob/goforum/model"
)

type PostStore struct {
	db *gorm.DB
}

func NewPostStore(db *gorm.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (as *PostStore) GetById(id uint) (*model.Post, error) {
	var m model.Post

	err := as.db.Where(id).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, err
}