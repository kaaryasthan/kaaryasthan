package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
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

type projectDS struct{}

func (ds *projectDS) Create(usr *user.User, prj *project.Project) error {
	prj.ID = 1
	return nil
}

func (ds *projectDS) Valid(prj *project.Project) error {
	return nil
}

func (ds *projectDS) Show(prj *project.Project) error {
	prj.ID = 1
	prj.Name = "somename"
	prj.Description = "Some description"
	prj.ItemTemplate = ""
	prj.Archived = false
	return nil
}

func (ds *projectDS) List(all bool) ([]project.Project, error) {
	return nil, nil
}

func TestProjectCreateHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewController(&userDS{}, &projectDS{})
	r.Handle("/api/v1/projects", negroni.New(
		negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.CreateHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "projects",
			"attributes": {
				"name": "somename",
				"description": "Some description"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(project.Project)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	respPayload := new(project.Project)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}
	reqPayload.ID = respPayload.ID

	if reqPayload.ID <= 0 {
		t.Errorf("ID is not 1 or above: %#v", reqPayload.ID)
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
