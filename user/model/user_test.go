package user

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/test"
)

func TestUserShow(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)

	usrDS := NewDatastore(DB)

	usr := &User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	usr2 := &User{Username: "jack"}
	if err := usrDS.Show(usr2); err == nil {
		t.Error("email is not yet verified")
	}

	DB.Exec("UPDATE users SET email_verified=true WHERE id=$1", usr.ID)

	usr3 := &User{Username: "jack"}
	if err := usrDS.Show(usr3); err != nil {
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
	if usr3.Active {
		t.Error("Wrong Active:", usr3.Active)
	}
	if !usr3.EmailVerified {
		t.Error("Wrong EmailVerified:", usr3.EmailVerified)
	}

	DB.Exec("UPDATE users SET active=true WHERE id=$1", usr.ID)

	usr4 := &User{Username: "jack"}
	if err := usrDS.Show(usr4); err != nil {
		t.Error(err)
	}

	if !usr4.Active {
		t.Error("Wrong Active:", usr4.Active)
	}

}
