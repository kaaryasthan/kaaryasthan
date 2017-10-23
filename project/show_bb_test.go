package project_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/urfave/negroni"
)

func TestProjectShowHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := project.NewController(&userDS{}, &projectDS{})
	r.Handle("/api/v1/projects/{name}", negroni.New(
		negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.ShowHandler)),
	)).Methods("GET")
	n.UseHandler(r)

	req, _ := http.NewRequest("GET", "/api/v1/projects/somename", nil)
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	if tr.Code != http.StatusOK {
		t.Error("User found with response:", tr.Code)
	}

	respPayload := new(project.Project)
	if err := jsonapi.UnmarshalPayload(tr.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}

	if respPayload.Name != "somename" {
		t.Error("Wrong Name:", respPayload.Name)
	}

	if respPayload.Description != "Some description" {
		t.Error("Wrong Description:", respPayload.Description)
	}

	if respPayload.ItemTemplate != "" {
		t.Error("Wrong item template:", respPayload.ItemTemplate)
	}

	if respPayload.Archived {
		t.Error("Archived:", respPayload.Archived)
	}
}
