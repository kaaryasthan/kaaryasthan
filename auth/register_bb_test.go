package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	controller "github.com/kaaryasthan/kaaryasthan/auth"
	auth "github.com/kaaryasthan/kaaryasthan/auth/model"
	correspondence "github.com/kaaryasthan/kaaryasthan/correspondence/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/urfave/negroni"
)

type userDS struct{}

func (ds *userDS) Create(usr *user.User) error {
	usr.ID = "1"
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

type authDS struct{}

func (ds *authDS) Login(obj *auth.Login) error {
	return nil
}

func (ds *authDS) VerifyEmail(obj *auth.Login) error {
	return nil
}

type correspondenceDS struct{}

func (ds *correspondenceDS) SendMail(eml *correspondence.Email) error {
	return nil
}

func TestUserRegisterHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewController(&userDS{}, &authDS{}, &correspondenceDS{})
	r.Handle("/api/v1/register", negroni.New(
		negroni.HandlerFunc(controller.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.RegisterHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "users",
			"attributes": {
				"username": "somename",
				"email": "somename@example.org",
				"password": "secret"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(user.User)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
	}

	respPayload := new(user.User)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}

	reqPayload.ID = respPayload.ID
	reqPayload.Password = ""

	if reqPayload.ID == "" {
		t.Fatalf("Login ID is empty")
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Fatalf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
