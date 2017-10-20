package item_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestItemShowHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")
	defer db.DB.Exec("DELETE FROM items")
	defer db.DB.Exec("DELETE FROM item_discussion_comment_search")

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

	itm := item.Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	err := itm.Create(usr)
	if err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != 1 {
		t.Fatalf("Data not inserted. Num: %#v", itm.Number)
	}

	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/items/"+strconv.Itoa(itm.Number), nil)
	req.Header.Set("Authorization", "Bearer "+tkn)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error("User found with response:", resp.Status)
	}
	respPayload := new(item.Item)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
		t.Error("Unable to unmarshal body:", err)
	}

	if respPayload.ID == 0 {
		t.Error("ID not set")
	}

	if respPayload.Description != "Some description" {
		t.Error("Wrong Description:", respPayload.Description)
	}
}
