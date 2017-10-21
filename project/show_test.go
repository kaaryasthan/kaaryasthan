package project

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestProjectShow(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")

	usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Fatal(err)
	}

	prj := Project{Name: "somename", Description: "Some description"}
	if err := prj.Create(usr); err != nil {
		t.Fatal(err)
	}
	if prj.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", prj.ID)
	}

	prj2 := Project{Name: "somename"}
	if err := prj2.Show(); err != nil {
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
