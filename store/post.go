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

	err := as.db.Where(id).Preload("User").First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, err
}

func (as *PostStore) List(offset, limit int) ([]model.Post, int, error) {
	var (
		posts []model.Post
		count int
	)

	as.db.Model(&posts).Count(&count)
	as.db.Offset(offset).
		Limit(limit).
		Order("created_at desc").Find(&posts)

	return posts, count, nil
}

func (as *PostStore) CreatePost(a *model.Post) error {
	tx := as.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where(a.ID).Preload("User").Find(&a).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (as *PostStore) UpdatePost(a *model.Post) error {
	tx := as.db.Begin()
	if err := tx.Model(a).Update(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where(a.ID).Preload("User").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (as *PostStore) DeletePost(a *model.Post) error {
	return as.db.Delete(a).Error
}

func (as *PostStore) AddComment(a *model.Post, c *model.Comment) error {
	err := as.db.Model(a).Association("Comments").Append(c).Error
	if err != nil {
		return err
	}

	return as.db.Where(c.ID).Preload("User").First(c).Error
}

func (as *PostStore) GetCommentsByPostId(id uint) ([]model.Comment, error) {
	var m model.Post
	err := as.db.Where(id).Preload("Comments").Preload("Comments.User").First(&m).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return m.Comments, nil
}

func (as *PostStore) GetCommentByID(id uint) (*model.Comment, error) {
	var m model.Comment
	if err := as.db.Where(id).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (as *PostStore) DeleteComment(c *model.Comment) error {
	return as.db.Delete(c).Error
}
