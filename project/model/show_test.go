package project

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestProjectShow(t *testing.T) {
	t.Parallel()
	DB := test.NewTestDB()
	defer test.ResetDB(DB)

	usrDS := user.NewDatastore(DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := NewDatastore(DB)
	prj := &Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Fatal(err)
	}
	if prj.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", prj.ID)
	}

	prj2 := &Project{Name: "somename"}
	if err := prjDS.Show(prj2); err != nil {
		t.Error("Project is valid", err)
	}

	if prj2.ID == 0 {
		t.Error("ID not set")
	}

	if prj2.Description != "Some description" {
		t.Error("Wrong Description", prj2.Description)
	}

	if prj2.Archived {
		t.Error("Wrong Archived", prj2.Archived)
	}
}
