package main

import (
	"github.com/shkryob/goforum/db"
	"github.com/shkryob/goforum/handler"
	"github.com/shkryob/goforum/router"
	"github.com/shkryob/goforum/store"
)

func main() {
	r := router.New()

	v1 := r.Group("/api")

	d := db.New()
	db.AutoMigrate(d)

	us := store.NewUserStore(d)
	as := store.NewPostStore(d)
	h := handler.NewHandler(us, as)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:1323"))
}
