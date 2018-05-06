// +build integration

package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth/model"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestIntegration(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)

	bi := search.NewBleveIndex(DB, conf)
	listener := bi.SubscribeAndCreateIndex()
	defer listener.Close()

	_, _, urt := route.Router(DB, bi)
	ts := httptest.NewServer(urt)
	defer ts.Close()

	tkn, usr := func() (string, *user.User) {

		usrDS := user.NewDatastore(DB)
		usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
		if err := usrDS.Create(usr); err != nil {
			t.Fatal(err)
		}

		DB.Exec("UPDATE users SET active=true, email_verified=true WHERE id=$1", usr.ID)
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
		return respPayload.Token, usr
	}()

	prjDS := project.NewDatastore(DB)
	prj := &project.Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Fatal(err)
	}

	itmDS := item.NewDatastore(DB, bi)
	itm := &item.Item{Title: "sometitle", Description: "Some description", ProjectID: strconv.Itoa(prj.ID)}
	if err := itmDS.Create(usr, itm); err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != "1" {
		t.Fatalf("Data not inserted. Num: %#v", itm.Number)
	}

	discDS := item.NewDiscussionDatastore(DB)
	disc := &item.Discussion{Body: "some discussion", ItemID: itm.ID}
	if err := discDS.Create(usr, disc); err != nil {
		t.Fatal(err)
	}
	if disc.ID == "" {
		t.Fatalf("Data not inserted. ID: %#v", disc.ID)
	}

	n := []byte(`{
			"data": {
				"type": "discussions",
				"attributes": {
					"body": "Some Body"
				}
			}
		}`)

	reqPayload := new(item.Discussion)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(n), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/items/1/relationships/discussions", bytes.NewReader(n))
	req.Header.Set("Authorization", "Bearer "+tkn)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := new(item.Discussion)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
		return
	}
	reqPayload.ItemID = 1

	reqPayload.ID = respPayload.ID

	if reqPayload.ID == "" {
		t.Error("Discussion ID is empty")
	}
	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}

	req2, _ := http.NewRequest("GET", ts.URL+"/api/v1/items?Query=sometitle", nil)
	req2.Header.Set("Authorization", "Bearer "+tkn)
	client2 := http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	buf := test.BufLog(t, resp2.Body, "Item:")
	if _, err := jsonapi.UnmarshalManyPayload(buf, reflect.TypeOf(new(item.Item))); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
		return
	}
}
