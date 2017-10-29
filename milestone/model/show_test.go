package milestone

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestMilestoneShow(t *testing.T) {
	t.Parallel()
	DB := test.NewTestDB()
	defer test.ResetDB(DB)

	usrDS := user.NewDatastore(DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := project.NewDatastore(DB)
	prj := &project.Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Error(err)
	}

	milDS := NewDatastore(DB)
	mil := &Milestone{Name: "somename", Description: "Some description", ProjectID: prj.ID}
	if err := milDS.Create(usr, mil); err != nil {
		t.Fatal(err)
	}
	if mil.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", mil.ID)
	}

	mil2 := &Milestone{Name: "somename", ProjectID: prj.ID}
	if err := milDS.Show(mil2); err != nil {
		t.Error("Milestone is valid", err)
	}

	if mil2.Description != "Some description" {
		t.Error("Wrong Description", mil2.Description)
	}

	if mil2.Items != nil {
		t.Error("Wrong Items", mil2.Items)
	}
}
