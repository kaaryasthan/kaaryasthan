package item

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestItemCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM items")
	s := Schema{Title: "sometitle", Description: "Some description"}
	id, err := s.Create()
	if err != nil {
		t.Error(err)
	}
	if id <= 0 {
		t.Errorf("Data not inserted. ID: %#v", id)
	}
}
