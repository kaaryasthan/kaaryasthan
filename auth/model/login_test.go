package auth

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestUserLogin(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)

	loginDS := NewDatastore(DB)

	login := &Login{Username: "jack", Password: "Secret@123"}
	if err := loginDS.Login(login); err == nil {
		t.Error("Login succeeded")
	}

	usrDS := user.NewDatastore(DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	if err := loginDS.Login(login); err == nil {
		t.Error("Login succeeded")
	}

	if _, err := DB.Exec("UPDATE users SET active=true"); err != nil {
		t.Fatal("Unable to updates users:", err)
	}
	if err := loginDS.Login(login); err == nil {
		t.Error("Login succeeded")
	}

	if _, err := DB.Exec("UPDATE users SET active=false, email_verified=true"); err != nil {
		t.Fatal("Unable to updates users:", err)
	}
	if err := loginDS.Login(login); err == nil {
		t.Error("Login succeeded")
	}

	if _, err := DB.Exec("UPDATE users SET active=true, email_verified=true"); err != nil {
		t.Fatal("Unable to updates users:", err)
	}
	if err := loginDS.Login(login); err != nil {
		t.Error("Login failed", err)
	}
	if login.ID == "" {
		t.Error("Login ID not set")
	}
}
