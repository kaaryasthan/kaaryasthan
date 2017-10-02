package discussion

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/item"
)

func TestDiscussionCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM items")
	defer db.DB.Exec("DELETE FROM discussions")
	is := item.Schema{Title: "sometitle", Description: "Some description"}
	id, err := is.Create()
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Fatalf("Item not inserted. ID: %#v", id)
	}

	s := Schema{Body: "some discussion", Item: id}
	id, err = s.Create()
	if err != nil {
		t.Error(err)
	}
	if id <= 0 {
		t.Errorf("Data not inserted. ID: %#v", id)
	}
}
