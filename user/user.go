package user

import (
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/baijum/logger"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/db"
	"golang.org/x/crypto/scrypt"
)

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
	PersonalNote  string `jsonapi:"attr,peronal_note,omitempty"`
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

// Show a user
func (obj *User) Show() error {
	err := db.DB.QueryRow(`SELECT id, name, email, user_role, active, email_verified, personal_note FROM "users"
		WHERE username=$1 AND email_verified=true AND deleted_at IS NULL`,
		obj.Username).Scan(&obj.ID, &obj.Name, &obj.Email, &obj.Role, &obj.Active, &obj.EmailVerified, &obj.PersonalNote)
	return err
}

func showUserHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)

	usr := User{ID: userID}
	if err := usr.Valid(); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't validate user: "+usr.ID, err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	obj := &User{Username: username}
	if err := obj.Show(); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't find user: "+username, err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		if logger.Level <= logger.DEBUG {
			log.Println("Couldn't unmarshal: ", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Register handlers
func Register(art, urt *mux.Router) {
	// art.HandleFunc("/api/v1/users", listUsersHandler).Methods("GET")
	art.HandleFunc("/api/v1/users/{username}", showUserHandler).Methods("GET")
	//art.HandleFunc("/api/v1/users/{username}", updateUserHandler).Methods("PATCH")
	//art.HandleFunc("/api/v1/users/{username}", deleteUserHandler).Methods("DELETE")
}
