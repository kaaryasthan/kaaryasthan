package item

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestItemList(t *testing.T) {
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

	for i := 0; i < 3; i++ {
		itm := &Item{Title: "found sometitle" + strconv.Itoa(i), Description: "Some awesome description" + strconv.Itoa(i), ProjectID: prj.ID,
			Labels: []string{"a", "c"}}
		err := itmDS.Create(usr, itm)
		if err != nil {
			t.Fatal(err)
		}
		if itm.ID <= 0 {
			t.Fatalf("Data not inserted. ID: %#v", itm.ID)
		}
		if itm.Number != i+1 {
			t.Fatalf("Data not inserted. Num: %#v", itm.Number)
		}
	}

	itm := &Item{Title: "found baiju sometitle", Description: "Some awesome description", ProjectID: prj.ID, Labels: []string{"a/c -_d", "b"}}
	err := itmDS.Create(usr, itm)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second) // FIXME: Any other approach possible?
	items, err := itmDS.List(`label:"a/c -_d"`, 0, 20)
	for _, i := range items {
		fmt.Println(i)
	}
	if err != nil {
		t.Error("Retrieving items failed.")
	}
	/*
		if items[0].ID != 1 {
			t.Error("Wrong ID", items[0].ID)
		}
		if items[0].Title != "found sometitle0" {
			t.Error("Wrong Title", items[0].Title)
		}
		if items[0].Description != "Some awesome description0" {
			t.Error("Wrong Title", items[0].Description)
		}
	*/
}
