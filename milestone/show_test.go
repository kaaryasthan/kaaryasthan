package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	authctrl "github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/milestone"
	"github.com/kaaryasthan/kaaryasthan/milestone/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/urfave/negroni"
)

func TestMilestoneShowHandler(t *testing.T) {
	t.Parallel()

	n := negroni.New()
	r := mux.NewRouter()
	c := controller.NewController(&userDS{}, &projectDS{}, &milestoneDS{})
	r.Handle("/api/v1/projects/{project}/milestones/{name}", negroni.New(
		negroni.HandlerFunc(authctrl.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(c.ShowHandler)),
	)).Methods("GET")
	n.UseHandler(r)

	req, _ := http.NewRequest("GET", "/api/v1/projects/p1/milestones/somename", nil)
	req.Header.Set("Authorization", test.NewBearerToken())
	tr := httptest.NewRecorder()
	n.ServeHTTP(tr, req)

	if tr.Code != http.StatusOK {
		t.Error("Milestone found with response:", tr.Code)
	}

	respPayload := new(milestone.Milestone)
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
}
