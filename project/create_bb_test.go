package project_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestProjectCreateHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")

	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()

	tkn := func() string {
		usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
		if err := usr.Create(); err != nil {
			t.Log("User creation failed", err)
			t.FailNow()
		}
		db.DB.Exec("UPDATE users SET active=true, email_verified=true WHERE id=$1", usr.ID)
		n := []byte(`{
			"data": {
				"type": "logins",
				"attributes": {
					"username": "jack",
					"password": "Secret@123"
				}
			}
		}`)

		reqPayload := new(auth.Login)
		if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
			t.Fatal("Unable to unmarshal input:", err)
		}

		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/login", bytes.NewReader(n))
		client := http.Client{}
		var err error
		var resp *http.Response
		if resp, err = client.Do(req); err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		respPayload := new(auth.Login)
		if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
			t.Fatal("Unable to unmarshal body:", err)
		}

		reqPayload.ID = respPayload.ID
		reqPayload.Token = respPayload.Token
		respPayload.Password = reqPayload.Password

		if reqPayload.ID == "" {
			t.Fatalf("Login ID is empty")
		}

		if !reflect.DeepEqual(reqPayload, respPayload) {
			t.Fatalf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
		}
		return respPayload.Token
	}()

	n := []byte(`{
		"data": {
			"type": "projects",
			"attributes": {
				"name": "somename",
				"description": "Some description"
			}
		}
	}`)
	reqPayload := new(project.Project)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/projects", bytes.NewReader(n))
	req.Header.Set("Authorization", "Bearer "+tkn)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := new(project.Project)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
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
