package item_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/db"
	. "github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestCommentCreateHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")
	defer db.DB.Exec("DELETE FROM items")
	defer db.DB.Exec("DELETE FROM item_discussion_comment_search")
	defer db.DB.Exec("DELETE FROM discussions")
	defer db.DB.Exec("DELETE FROM comments")

	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()

	tkn, usr := func() (string, user.User) {
		usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
		err := usr.Create()
		if err != nil {
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

		if reqPayload.ID == "" {
			t.Fatalf("Login ID is empty")
		}

		if !reflect.DeepEqual(reqPayload, respPayload) {
			t.Fatalf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
		}
		return respPayload.Token, usr
	}()

	prj := project.Project{Name: "somename", Description: "Some description"}
	if err := prj.Create(usr); err != nil {
		t.Fatal(err)
	}

	itm := Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	if err := itm.Create(usr); err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Num != 1 {
		t.Fatalf("Data not inserted. Num: %#v", itm.Num)
	}

	disc := Discussion{Body: "some discussion", ItemID: itm.ID}
	if err := disc.Create(usr); err != nil {
		t.Fatal(err)
	}
	if disc.ID == "" {
		t.Fatalf("Data not inserted. ID: %#v", disc.ID)
	}

	n := []byte(fmt.Sprintf(`{
		"data": {
			"type": "comments",
			"attributes": {
				"body": "Some Body",
				"discussion_id": "%s"
			}
		}
	}`, disc.ID))

	reqPayload := new(Comment)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/comments", bytes.NewReader(n))
	req.Header.Set("Authorization", "Bearer "+tkn)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := new(Comment)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
		return
	}

	reqPayload.ID = respPayload.ID

	if reqPayload.ID == "" {
		t.Error("Comment ID is empty")
	}
	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}

	/*
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
	*/
}
