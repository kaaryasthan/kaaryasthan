package comment_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	. "github.com/kaaryasthan/kaaryasthan/comment"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/route"
)

func TestCommentCreateHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM items")
	defer db.DB.Exec("DELETE FROM comments")
	_, art, _ := route.Router()
	ts := httptest.NewServer(art)
	defer ts.Close()
	n1 := []byte(`{
  "data": {
    "type": "items",
    "attributes": {
      "title": "Some Title",
      "description": "Some description"
    }
  }
}`)
	req1, _ := http.NewRequest("POST", ts.URL+"/api/v1/items", bytes.NewReader(n1))
	client1 := http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		log.Println("Unable to make request: ", err)
		return
	}
	defer resp1.Body.Close()

	respItemPayload := make(map[string]item.Data)
	itemDecoder := json.NewDecoder(resp1.Body)
	itemDecoder.Decode(&respItemPayload)
	itemID, _ := strconv.Atoi(respItemPayload["data"].ID)

	n := []byte(fmt.Sprintf(`{
  "data": {
    "type": "comments",
    "attributes": {
      "body": "Some comment"
    },
    "relationships": {
	    "items" : {
		    "data": {
			    "type": "items",
			    "id": "%d"
		    }
	    }
    }
  }
}`, itemID))
	reqPayload := make(map[string]Data)
	decoder1 := json.NewDecoder(bytes.NewReader(n))
	err = decoder1.Decode(&reqPayload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/items/"+fmt.Sprintf("%d", itemID)+"/comments", bytes.NewReader(n))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := make(map[string]Data)
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
