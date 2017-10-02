package milestone_test

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

func TestMilestoneCreateHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM milestones")
	_, art, _ := route.Router()
	ts := httptest.NewServer(art)
	defer ts.Close()
	n := []byte(`{
  "data": {
    "type": "milestones",
    "attributes": {
      "name": "somename",
      "description": "Some description"
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

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/milestones", bytes.NewReader(n))
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
