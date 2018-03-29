package auth

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/scrypt"
)

// VerifyEmail verify user
func (ds *Datastore) VerifyEmail(obj *Login) error {
	var originalPassword, salt []byte
	var originalEVC string
	err := ds.db.QueryRow(`SELECT id, password, salt, email_verification_code FROM "users"
		WHERE username=$1 AND email_verified=false`,
		obj.Username).Scan(&obj.ID, &originalPassword, &salt, &originalEVC)
	if err != nil {
		return err
	}

	if obj.ID == "" {
		return errors.New("User doesn't exist")
	}

	newPassword, err := scrypt.Key([]byte(obj.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}
	if !bytes.Equal(newPassword, originalPassword) {
		return errors.New("Wrong username or password")
	}
	if originalEVC != obj.EmailVerificationCode {
		return errors.New("Wrong email verification code")
	}
	_, err = ds.db.Exec("UPDATE users SET active=true, email_verified=true")
	return err
}
