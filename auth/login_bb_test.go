package auth_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/route"
)

func TestUserLoginHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM members")
	s2 := Schema{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	err := s2.Register()
	if err != nil {
		t.Log("Registration failed", err)
		t.FailNow()
	}
	db.DB.Exec("UPDATE members SET active=true, email_verified=true")
	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "members",
    "attributes": {
      "username": "jack",
      "password": "Secret@123"
    }
  }
}`)
	reqPayload := make(map[string]jsonapi.Data)
	decoder1 := json.NewDecoder(bytes.NewReader(n))
	err = decoder1.Decode(&reqPayload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/login", bytes.NewReader(n))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respPayload := make(map[string]jsonapi.Data)
	decoder2 := json.NewDecoder(resp.Body)
	err = decoder2.Decode(&respPayload)
	if err != nil {
		t.Fatal("Unable to decode body: ", err)
	}
	respData := respPayload["data"]
	reqData := reqPayload["data"]
	reqData.ID = respData.ID
	delete(reqData.Attributes, "password")
	delete(reqData.Attributes, "name")
	delete(reqData.Attributes, "email")
	if !reflect.DeepEqual(reqData, respData) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqData, respData)
	}
}
