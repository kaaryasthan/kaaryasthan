package route

import (
	"os"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/test"
)

func TestMain(m *testing.M) {
	dbname := test.NewTestDB()
	code := m.Run()
	test.ResetDB(dbname)
	os.Exit(code)
}
