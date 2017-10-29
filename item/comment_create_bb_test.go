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

type commentDS struct{}

func (ds *commentDS) Create(usr *user.User, cmt *item.Comment) error {
	cmt.ID = "8291a250-3e6f-4743-9153-64876b2c8fde"
	return nil
}

func TestCommentCreateHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewCommentController(&userDS{}, &discussionDS{}, &commentDS{})
	r.Handle("/api/v1/comments", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.CreateCommentHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "comments",
			"attributes": {
				"body": "Some body",
				"discussion_id": "5bb5faaf-48f5-4ff7-95dd-9bf6088956f3"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/comments", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(item.Comment)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	respPayload := new(item.Comment)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}
	reqPayload.ID = respPayload.ID

	if reqPayload.ID != "8291a250-3e6f-4743-9153-64876b2c8fde" {
		t.Errorf("ID is not '8291a250-3e6f-4743-9153-64876b2c8fde' got: %#v", reqPayload.ID)
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
