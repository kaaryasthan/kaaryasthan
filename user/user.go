package user

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/baijum/logger"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/scrypt"
)

// Repository helps to manage users
type Repository interface {
	Create(usr *User) error
	Valid(usr *User) error
	Show(usr *User) error
}

// Controller holds DB
type Controller struct {
	ds Repository
}

// User represents a user
type User struct {
	ID            string `jsonapi:"primary,users"`
	Username      string `jsonapi:"attr,username"`
	Name          string `jsonapi:"attr,name"`
	Email         string `jsonapi:"attr,email"`
	Role          string `jsonapi:"attr,role"`
	Active        bool   `jsonapi:"attr,active"`
	EmailVerified bool   `jsonapi:"attr,email_verified"`
	Password      string `jsonapi:"attr,password,omitempty"`
	PersonalNote  string `jsonapi:"attr,personal_note,omitempty"`
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}

// Create a new user
func (ds *Datastore) Create(usr *User) error {
	salt := randomSalt()
	password, err := scrypt.Key([]byte(usr.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}
	err = ds.db.QueryRow(`INSERT INTO "users"
		(username, name, email, password, salt)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		usr.Username,
		usr.Name,
		usr.Email,
		password,
		salt).Scan(&usr.ID)
	return err
}

// Valid checks the validity of the user
func (ds *Datastore) Valid(usr *User) error {
	var count int
	err := ds.db.QueryRow(`SELECT count(1) FROM "users"
		WHERE id=$1 AND active=true AND email_verified=true AND deleted_at IS NULL`,
		usr.ID).Scan(&count)
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

// Show a user
func (ds *Datastore) Show(usr *User) error {
	err := ds.db.QueryRow(`SELECT id, name, email, user_role, active, email_verified, personal_note FROM "users"
		WHERE username=$1 AND email_verified=true AND deleted_at IS NULL`,
		usr.Username).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Role, &usr.Active, &usr.EmailVerified, &usr.PersonalNote)
	return err
}

// ShowUserHandler get one user
func (c *Controller) ShowUserHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)

	usr := &User{ID: userID}
	if err := c.ds.Valid(usr); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't validate user: "+usr.ID, err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	usr = &User{Username: username}
	if err := c.ds.Show(usr); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't find user: "+username, err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := jsonapi.MarshalPayload(w, usr); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't unmarshal: ", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// NewController constructs a controller
func NewController(repo Repository) *Controller {
	return &Controller{ds: repo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(NewDatastore(db))

	// art.HandleFunc("/api/v1/users", listUsersHandler).Methods("GET")
	art.HandleFunc("/api/v1/users/{username}", c.ShowUserHandler).Methods("GET")
	//art.HandleFunc("/api/v1/users/{username}", updateUserHandler).Methods("PATCH")
	//art.HandleFunc("/api/v1/users/{username}", deleteUserHandler).Methods("DELETE")
}
