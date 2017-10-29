package auth

import (
	"database/sql"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/config"
	user "github.com/kaaryasthan/kaaryasthan/user/model"
)

// JwtMiddleware is middleware to handle all request
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(config.Config.TokenSecretKey)
		return secretKey, nil
	},
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
})

// Repository helps to manage login
type Repository interface {
	Login(obj *Login) error
}

// Controller holds DB
type Controller struct {
	ds  Repository
	uds user.Repository
}

// Login represents a logged in user session
type Login struct {
	ID       string `jsonapi:"primary,logins"`
	Username string `jsonapi:"attr,username"`
	Password string `jsonapi:"attr,password,omitempty"`
	Token    string `jsonapi:"attr,token,omitempty"`
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}

// NewController constructs a controller
func NewController(userRepo user.Repository, repo Repository) *Controller {
	return &Controller{ds: repo, uds: userRepo}
}

// Register handlers
func Register(urt *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), NewDatastore(db))
	urt.HandleFunc("/api/v1/register", c.RegisterHandler).Methods("POST")
	urt.HandleFunc("/api/v1/login", c.LoginHandler).Methods("POST")
}
