package db

import (
	"github.com/labstack/gommon/log"

	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/shkryob/goforum/model"
)

func New() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./goforum.db")
	if err != nil {
		log.Fatal("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./goforum_test.db")
	if err != nil {
		log.Fatal("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(false)
	return db
}

func DropTestDB() error {
	if err := os.Remove("./goforum_test.db"); err != nil {
		return err
	}
	return nil
}

//TODO: err check
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)
}
