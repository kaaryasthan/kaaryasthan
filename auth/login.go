package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/scrypt"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
)

func (obj *Schema) login() error {
	var originalPassword, salt []byte
	err := db.DB.QueryRow(`SELECT password, salt FROM "members"
		WHERE username=$1 AND active=true AND email_verified=true`,
		obj.Username).Scan(&originalPassword, &salt)
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
	payload := make(map[string]jsonapi.Data)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}
	s := New(payload["data"])
	err = s.login()
	if err != nil {
		log.Println("Login verification failed: ", err)
		return
	}
	tmpData := payload["data"]
	tmpData.ID = s.Username
	delete(tmpData.Attributes, "password")
	delete(tmpData.Attributes, "name")
	delete(tmpData.Attributes, "email")
	payload["data"] = tmpData
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Unable marshal data: ", err)
		return
	}
	w.Write(b)
}
