package user

import (
	"crypto/rand"
	"errors"
	"io"
	"log"

	"github.com/kaaryasthan/kaaryasthan/db"
	"golang.org/x/crypto/scrypt"
)

// User represents a user
type User struct {
	ID       string `jsonapi:"primary,users"`
	Username string `jsonapi:"attr,username"`
	Name     string `jsonapi:"attr,name"`
	Email    string `jsonapi:"attr,email"`
	Role     string `jsonapi:"attr,role"`
	Password string `jsonapi:"attr,password,omitempty"`
}

// Create a new user
func (obj *User) Create() error {
	salt := randomSalt()
	password, err := scrypt.Key([]byte(obj.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}
	err = db.DB.QueryRow(`INSERT INTO "users"
		(username, name, email, password, salt)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		obj.Username,
		obj.Name,
		obj.Email,
		password,
		salt).Scan(&obj.ID)
	return err
}

// Valid checks the validity of the user
func (obj *User) Valid() error {
	var count int
	err := db.DB.QueryRow(`SELECT count(1) FROM "users"
		WHERE id=$1 AND active=true AND email_verified=true AND deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid user")
	}
	return nil
}

func randomSalt() []byte {
	s := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, s)
	if err != nil {
		log.Println(err)
	}
	return s
}
