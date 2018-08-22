package test

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kaaryasthan/kaaryasthan/config"
)

// NewBearerToken generates a token with 1 minute validity
func NewBearerToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1b54ad59-1624-4f94-8ac9-fb90731620eb",
		"exp": time.Now().Add(time.Minute).Unix(),
	})

	tkn, err := token.SignedString([]byte(config.Config.TokenSecretKey))
	if err != nil {
		log.Println("Error sigining", err)
		return ""
	}
	return "Bearer " + tkn
}
