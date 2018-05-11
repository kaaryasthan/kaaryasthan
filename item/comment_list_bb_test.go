package controller_test

import (
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
	"github.com/urfave/negroni"
)

func TestCommentListHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewCommentController(&userDS{}, &itemDS{}, &commentDS{})
	r.Handle("/api/v1/items/{number:[1-9]\\d*}/comments", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.ListCommentHandler)),
	)).Methods("GET")
	n.UseHandler(r)

	req, _ := http.NewRequest("GET", "/api/v1/items/1/comments", nil)
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	if tr.Code != http.StatusOK {
		t.Error("User found with response:", tr.Code)
	}

	buf := test.BufLog(t, tr.Body, "Comment:")

	respPayload, err := jsonapi.UnmarshalManyPayload(buf, reflect.TypeOf(new(item.Comment)))
	if err != nil {
		t.Errorf("Unable to unmarshal body: %+v", err)
		return
	}

	comments := make([]*item.Comment, len(respPayload))
	for i, v := range respPayload {
		comments[i] = v.(*item.Comment)
	}

	if comments[0].Body != "Some body" {
		t.Error("Wrong Body:", comments[0].Body)
	}
}
