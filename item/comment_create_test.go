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
	item "github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/urfave/negroni"
)

type commentDS struct{}

func (ds *commentDS) Create(usr *user.User, cmt *item.Comment) error {
	cmt.ID = "5bb5faaf-48f5-4ff7-95dd-9bf6088956f3"
	return nil
}

func (ds *commentDS) Valid(cmt *item.Comment) error {
	return nil
}

func (ds *commentDS) List(itmID int) ([]*item.Comment, error) {
	dl := []*item.Comment{
		{ID: "1", Body: "Some body", ItemID: 1},
		{ID: "2", Body: "Some body 2", ItemID: 1},
	}
	return dl, nil
}

func TestCommentCreateHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewCommentController(&userDS{}, &itemDS{}, &commentDS{})
	r.Handle("/api/v1/items/{number:[1-9]\\d*}/relationships/comments", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.CreateCommentHandler)),
	)).Methods("POST")
	n.UseHandler(r)

	d := []byte(`{
		"data": {
			"type": "comments",
			"attributes": {
				"body": "Some body"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/items/1/relationships/comments", bytes.NewReader(d))
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	reqPayload := new(item.Comment)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(d), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}
	reqPayload.ItemID = 1

	respPayload := new(item.Comment)
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
