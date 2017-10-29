package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	authctrl "github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/urfave/negroni"
)

type discussionDS struct{}

func (ds *discussionDS) Create(usr *user.User, disc *item.Discussion) error {
	disc.ID = "5bb5faaf-48f5-4ff7-95dd-9bf6088956f3"
	return nil
}

func (ds *discussionDS) Valid(itm *item.Discussion) error {
	return nil
}

func TestDiscussionCreateHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewDiscussionController(&userDS{}, &itemDS{}, &discussionDS{})
	r.Handle("/api/v1/discussions", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.CreateDiscussionHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "discussions",
			"attributes": {
				"body": "Some body",
				"item_id": 1
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/discussions", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(item.Discussion)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	respPayload := new(item.Discussion)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}
	reqPayload.ID = respPayload.ID

	if reqPayload.ID != "5bb5faaf-48f5-4ff7-95dd-9bf6088956f3" {
		t.Errorf("ID is not '5bb5faaf-48f5-4ff7-95dd-9bf6088956f3' got: %#v", reqPayload.ID)
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
