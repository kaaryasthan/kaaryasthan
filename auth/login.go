package controller

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth/model"
	"github.com/kaaryasthan/kaaryasthan/config"
)

// LoginHandler login user
func (c *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(auth.Login)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.ds.Login(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": obj.ID,
		"exp": time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	secretKey := []byte(config.Config.TokenSecretKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj.Token = tokenString
	obj.Password = ""
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
