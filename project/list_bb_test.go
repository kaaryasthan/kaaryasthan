package controller_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	authctrl "github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/urfave/negroni"
)

func TestProjectListHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewController(&userDS{}, &projectDS{})
	r.Handle("/api/v1/projects", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.ListHandler)),
	)).Methods("GET")
	n.UseHandler(r)

	req, _ := http.NewRequest("GET", "/api/v1/projects", nil)
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	if tr.Code != http.StatusOK {
		t.Error("User found with response:", tr.Code)
	}

	respPayload, err := jsonapi.UnmarshalManyPayload(tr.Body, reflect.TypeOf(new(project.Project)))
	if err != nil {
		t.Error("Unable to unmarshal body:", err)
		return
	}

	projects := make([]*project.Project, len(respPayload))
	for i, v := range respPayload {
		projects[i] = v.(*project.Project)
	}

	if projects[0].Name != "somename1" {
		t.Error("Wrong Name:", projects[0].Name)
	}

	if projects[0].Description != "Some description 1" {
		t.Error("Wrong Description:", projects[0].Description)
	}

	if projects[0].ItemTemplate != "" {
		t.Error("Wrong item template:", projects[0].ItemTemplate)
	}

	if projects[0].Archived {
		t.Error("Archived:", projects[0].Archived)
	}

	if projects[1].Name != "somename2" {
		t.Error("Wrong Name:", projects[1].Name)
	}

	if projects[1].Description != "Some description 2" {
		t.Error("Wrong Description:", projects[1].Description)
	}

	if projects[1].ItemTemplate != "" {
		t.Error("Wrong item template:", projects[1].ItemTemplate)
	}

	if projects[1].Archived {
		t.Error("Archived:", projects[1].Archived)
	}

}
