package controller

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth/model"
	"github.com/kaaryasthan/kaaryasthan/config"
)

// EmailVerificationHandler register user
func (c *Controller) EmailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	obj := new(auth.Login)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.ds.VerifyEmail(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": obj.ID,
		"exp": time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	secretKey := []byte(config.Config.TokenSecretKey)
	tokenString, _ := token.SignedString(secretKey)
	obj.Token = tokenString
	obj.Password = ""
	obj.EmailVerificationCode = ""
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
