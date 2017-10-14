package auth

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
)

func (obj *Login) login() error {
	var originalPassword, salt []byte
	err := db.DB.QueryRow(`SELECT id, password, salt FROM "users"
		WHERE username=$1 AND active=true AND email_verified=true`,
		obj.Username).Scan(&obj.ID, &originalPassword, &salt)
	if err != nil {
		return err
	}

	newPassword, err := scrypt.Key([]byte(obj.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}
	if !bytes.Equal(newPassword, originalPassword) {
		return errors.New("Wrong username or password")
	}
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(Login)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj.login()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": obj.ID,
		"exp": time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	tokenString, _ := token.SignedString(secretKey)
	obj.Token = tokenString
	obj.Password = ""
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
