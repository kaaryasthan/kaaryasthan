package auth

import "database/sql"

// Repository helps to manage login
type Repository interface {
	Login(obj *Login) error
	VerifyEmail(obj *Login) error
}

// Login represents a logged in user session
type Login struct {
	ID                    string `jsonapi:"primary,logins"`
	Username              string `jsonapi:"attr,username"`
	Password              string `jsonapi:"attr,password,omitempty"`
	Token                 string `jsonapi:"attr,token,omitempty"`
	EmailVerificationCode string `jsonapi:"attr,email_verification_code,omitempty"`
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}
