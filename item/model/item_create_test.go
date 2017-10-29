package item

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestItemCreate(t *testing.T) {
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
		t.Fatal(err)
	}

	itmDS := NewDatastore(DB)
	itm := &Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	err := itmDS.Create(usr, itm)
	if err != nil {
		t.Error(err)
	}
	if itm.ID <= 0 {
		t.Errorf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != 1 {
		t.Errorf("Data not inserted. Num: %#v", itm.Number)
	}
}
