package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth/model"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestEmailVerificationHandler(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)
	bi := search.NewBleveIndex(DB, conf)

	usrDS := user.NewDatastore(DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	DB.Exec("UPDATE users SET email_verification_code='4208f094-b50d-4296-bdcc-166c76672c95'")
	_, _, urt := route.Router(DB, bi)
	ts := httptest.NewServer(urt)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "logins",
    "attributes": {
      "username": "jack",
      "password": "Secret@123",
      "email_verification_code": "4208f094-b50d-4296-bdcc-166c76672c95"
    }
  }
}`)

	reqPayload := new(auth.Login)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/emailverification", bytes.NewReader(n))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
	respPayload.EmailVerificationCode = reqPayload.EmailVerificationCode

	if reqPayload.ID == "" {
		t.Fatalf("Login ID is empty")
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Fatalf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
