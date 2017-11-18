package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestItemListHandler(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)
	bi := search.NewBleveIndex(DB, conf)

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
	itm := &item.Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	err := itmDS.Create(usr, itm)
	if err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != 1 {
		t.Fatalf("Data not inserted. Num: %#v", itm.Number)
	}

	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/items?Query=description", nil)
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
}
