package user_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestUserShowHandler(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")

	_, _, urt := route.Router()
	ts := httptest.NewServer(urt)
	defer ts.Close()

	tkn := func() string {
		usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
		if err := usr.Create(); err != nil {
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
		return respPayload.Token
	}()

	usr2 := user.User{Username: "jill", Name: "Jill Wilber", Email: "jill@example.com", Password: "Secret@123"}
	err := usr2.Create()
	if err != nil {
		t.Log("User creation failed", err)
		t.FailNow()
	}
	t.Run("email verified user", func(t *testing.T) {
		db.DB.Exec("UPDATE users SET active=true, email_verified=true WHERE id=$1", usr2.ID)
		req, _ := http.NewRequest("GET", ts.URL+"/api/v1/users/"+usr2.Username, nil)
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
		respPayload := new(user.User)
		if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
			t.Error("Unable to unmarshal body:", err)
			return
		}
		if respPayload.Username != "jill" {
			t.Error("Wrong Username:", respPayload.Username)
		}
		if respPayload.Name != "Jill Wilber" {
			t.Error("Wrong Name:", respPayload.Name)
		}
		if respPayload.Email != "jill@example.com" {
			t.Error("Wrong Eamil:", respPayload.Email)
		}
		if respPayload.Role != "member" {
			t.Error("Wrong Role:", respPayload.Role)
		}
		if !respPayload.Active {
			t.Error("Wrong Active:", respPayload.Active)
		}
		if !respPayload.EmailVerified {
			t.Error("Wrong EmailVerified:", respPayload.EmailVerified)
		}
	})
	t.Run("email not verified user", func(t *testing.T) {
		db.DB.Exec("UPDATE users SET active=true, email_verified=false WHERE id=$1", usr2.ID)
		req, _ := http.NewRequest("GET", ts.URL+"/api/v1/users/"+usr2.Username, nil)
		req.Header.Set("Authorization", "Bearer "+tkn)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Error("User found with response:", resp.Status)
		}
	})
}
