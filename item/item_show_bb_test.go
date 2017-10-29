package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	authctrl "github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/urfave/negroni"
)

func TestItemShowHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewItemController(&userDS{}, &projectDS{}, &itemDS{})
	r.Handle("/api/v1/items/{number:[1-9]\\d*}", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.ShowItemHandler)),
	)).Methods("GET")
	n.UseHandler(r)

	req, _ := http.NewRequest("GET", "/api/v1/items/1", nil)
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	if tr.Code != http.StatusOK {
		t.Error("User found with response:", tr.Code)
	}

	respPayload := new(item.Item)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}

	if respPayload.Title != "Some Title" {
		t.Error("Wrong Name:", respPayload.Title)
	}

	if respPayload.Description != "Some description" {
		t.Error("Wrong Description:", respPayload.Description)
	}
}
