package user_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user"
	"github.com/urfave/negroni"
)

type userDS struct{}

func (ds *userDS) Create(usr *user.User) error {
	return nil
}

func (ds *userDS) Valid(usr *user.User) error {
	return nil
}

func (ds *userDS) Show(usr *user.User) error {
	usr.Name = "Jack Wilber"
	usr.Email = "jack@example.com"
	usr.Role = "member"
	usr.Active = true
	usr.EmailVerified = true
	return nil
}

type invalidUserDS struct {
	userDS
}

func (ds *invalidUserDS) Valid(usr *user.User) error {
	return errors.New("Invalid user")
}

func TestUserShowHandler(t *testing.T) {
	t.Parallel()

	t.Run("valid user", func(t *testing.T) {
		t.Parallel()
		n := negroni.New()
		r := mux.NewRouter()
		c := user.NewController(&userDS{})
		r.Handle("/api/v1/users/{username}", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(http.HandlerFunc(c.ShowUserHandler)),
		)).Methods("GET")
		n.UseHandler(r)

		req, _ := http.NewRequest("GET", "/api/v1/users/jack", nil)
		req.Header.Set("Authorization", test.NewBearerToken())
		tr := httptest.NewRecorder()
		n.ServeHTTP(tr, req)

		if tr.Code != http.StatusOK {
			t.Error("User found with response:", tr.Code)
		}

		respPayload := new(user.User)
		if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
			t.Error("Unable to unmarshal body:", err)
			return
		}
		if respPayload.Username != "jack" {
			t.Error("Wrong Username:", respPayload.Username)
		}
		if respPayload.Name != "Jack Wilber" {
			t.Error("Wrong Name:", respPayload.Name)
		}
		if respPayload.Email != "jack@example.com" {
			t.Error("Wrong Eamil:", respPayload.Email)
		}
		if respPayload.Role != "member" {
			t.Error("Wrong Role:", respPayload.Role)
		}
		if !respPayload.Active {
			t.Error("Wrong Active:", respPayload.Active)
		}
		if !respPayload.EmailVerified {
			t.Error("Wrong EmailVerified:", respPayload.EmailVerified)
		}
	})

	t.Run("invalid user", func(t *testing.T) {
		t.Parallel()
		n := negroni.New()
		r := mux.NewRouter()
		c := user.NewController(&invalidUserDS{})
		r.Handle("/api/v1/users/{username}", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(http.HandlerFunc(c.ShowUserHandler)),
		)).Methods("GET")
		n.UseHandler(r)

		req, _ := http.NewRequest("GET", "/api/v1/users/jack", nil)
		req.Header.Set("Authorization", test.NewBearerToken())
		tr := httptest.NewRecorder()
		n.ServeHTTP(tr, req)

		if tr.Code != http.StatusUnauthorized {
			t.Error("User found with response:", tr.Code)
		}
	})
}
