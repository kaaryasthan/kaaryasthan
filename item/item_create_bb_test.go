package item_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/project"
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
	return nil
}

func (ds *projectDS) List(all bool) ([]project.Project, error) {
	return nil, nil
}

type itemDS struct{}

func (ds *itemDS) Create(usr *user.User, itm *item.Item) error {
	itm.ID = 1
	return nil
}

func (ds *itemDS) List(query string) ([]item.Item, error) {
	return nil, nil
}

func (ds *itemDS) Show(itm *item.Item) error {
	itm.Title = "Some Title"
	itm.Description = "Some description"
	return nil
}

func (ds *itemDS) Valid(itm *item.Item) error {
	return nil
}

func TestItemCreateHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := item.NewController(&userDS{}, &projectDS{}, &itemDS{})
	r.Handle("/api/v1/items", negroni.New(
		negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.CreateItemHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "items",
			"attributes": {
				"title": "Some Title",
				"description": "Some description",
				"project_id": 1
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/items", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(item.Item)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	respPayload := new(item.Item)
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
