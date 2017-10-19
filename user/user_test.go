package user

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestUserShow(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")

	usr := User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Fatal(err)
	}

	usr2 := User{Username: "jack"}
	if err := usr2.Show(); err == nil {
		t.Error("email is not yet verified")
	}

	db.DB.Exec("UPDATE users SET email_verified=true WHERE id=$1", usr.ID)

	usr3 := User{Username: "jack"}
	if err := usr3.Show(); err != nil {
		t.Error(err)
	}
	if usr3.Username != "jack" {
		t.Error("Wrong Username:", usr3.Username)
	}
	if usr3.Name != "Jack Wilber" {
		t.Error("Wrong Name:", usr3.Name)
	}
	if usr3.Email != "jack@example.com" {
		t.Error("Wrong Eamil:", usr3.Email)
	}
	if usr3.Role != "member" {
		t.Error("Wrong Role:", usr3.Role)
	}
	if usr3.Active != false {
		t.Error("Wrong Active:", usr3.Active)
	}
	if usr3.EmailVerified != true {
		t.Error("Wrong EmailVerified:", usr3.EmailVerified)
	}

	db.DB.Exec("UPDATE users SET active=true WHERE id=$1", usr.ID)

	usr4 := User{Username: "jack"}
	if err := usr4.Show(); err != nil {
		t.Error(err)
	}

	if usr4.Active != true {
		t.Error("Wrong Active:", usr4.Active)
	}

}
