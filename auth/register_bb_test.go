package auth_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/route"
)

func TestUserRegisterHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM members")
	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "members",
    "attributes": {
      "username": "jack",
      "name": "Jack Wilber",
      "email": "jack@example.com",
      "password": "Secret@123"
    }
  }
}`)
	reqPayload := make(map[string]jsonapi.Data)
	decoder1 := json.NewDecoder(bytes.NewReader(n))
	err := decoder1.Decode(&reqPayload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/register", bytes.NewReader(n))
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
	id, err := strconv.Atoi(reqData.ID)
	if err != nil {
		t.Fatal("Wrong ID", err)
	}
	if id <= 0 {
		t.Errorf("ID is not 1 or above: %#v", id)
	}
	if !reflect.DeepEqual(reqData, respData) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqData, respData)
	}
}
