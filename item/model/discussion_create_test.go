package item

import (
	"strconv"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestCommentCreate(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)
	bi := search.NewBleveIndex(DB, conf)
	listener := bi.SubscribeAndCreateIndex()
	defer listener.Close()

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

	itmDS := NewDatastore(DB, bi)
	itm := &Item{Title: "sometitle", Description: "Some description", ProjectID: strconv.Itoa(prj.ID)}
	if err := itmDS.Create(usr, itm); err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != "1" {
		t.Fatalf("Data not inserted. Num: %#v", itm.Number)
	}

	cmtDS := NewCommentDatastore(DB)
	cmt := &Comment{Body: "some comment", ItemID: itm.ID}
	if err := cmtDS.Create(usr, cmt); err != nil {
		t.Error(err)
	}
	if cmt.ID == "" {
		t.Errorf("Data not inserted. ID: %#v", cmt.ID)
	}
}
