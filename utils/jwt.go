package utils

import (
	"github.com/labstack/gommon/log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTSecret = []byte("!!SECRET!!")

func GenerateJWT(id uint) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
