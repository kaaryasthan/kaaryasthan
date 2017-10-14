package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestUserRegisterHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "users",
    "attributes": {
      "username": "jack",
      "name": "Jack Wilber",
      "email": "jack@example.com",
      "role": "admin",
      "password": "Secret@123"
    }
  }
}`)

	reqPayload := new(user.User)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/register", bytes.NewReader(n))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := new(user.User)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
	}

	reqPayload.ID = respPayload.ID
	reqPayload.Password = ""

	if reqPayload.ID == "" {
		t.Fatalf("Login ID is empty")
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Fatalf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
