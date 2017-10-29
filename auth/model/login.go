package auth

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/scrypt"
)

// Login verify user
func (ds *Datastore) Login(obj *Login) error {
	var originalPassword, salt []byte
	err := ds.db.QueryRow(`SELECT id, password, salt FROM "users"
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
