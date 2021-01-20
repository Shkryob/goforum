package handler

import (
	"fmt"
	"os"
	"testing"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/db"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/router"
	"github.com/shkryob/goforum/store"
)

var (
	d  *gorm.DB
	us UserStore
	as PostStore
	h  *Handler
	e  *echo.Echo
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func authHeader(token string) string {
	return "Token " + token
}

func setup() {
	d = db.TestDB()
	db.AutoMigrate(d)
	us = store.NewUserStore(d)
	as = store.NewPostStore(d)
	h = NewHandler(us, as)
	e = router.New()
	loadFixtures()
}

func tearDown() {
	_ = d.Close()
	if err := db.DropTestDB(); err != nil {
		fmt.Println(err)
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	u1 := model.User{
		Username: "user1",
		Email:    "user1@test.io",
	}
	u1.Password, _ = u1.HashPassword("secret")
	if err := us.Create(&u1); err != nil {
		return err
	}

	u2 := model.User{
		Username: "user2",
		Email:    "user2@test.io",
	}
	u2.Password, _ = u2.HashPassword("secret")
	if err := us.Create(&u2); err != nil {
		return err
	}

	a := model.Post{
		Title:  "post1 title",
		Body:   "post1 body",
		UserID: 1,
	}
	as.CreatePost(&a)
	as.AddComment(&a, &model.Comment{
		Body:   "post1 comment1",
		PostID: 1,
		UserID: 1,
	})

	a2 := model.Post{
		Title:  "post2 title",
		Body:   "post2 body",
		UserID: 2,
	}
	as.CreatePost(&a2)
	as.AddComment(&a2, &model.Comment{
		Body:   "post2 comment1 by user1",
		PostID: 2,
		UserID: 1,
	})

	return nil
}
