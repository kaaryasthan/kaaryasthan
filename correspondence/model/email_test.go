// +build external

package correspondence

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/test"
)

func TestSendMail(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)

	ds := NewDatastore(DB)
	eml := &Email{ToName: "Baiju Muthukadan",
		ToEmail:     "baiju.m.mail@gmail.com",
		Sender:      "noreply@kaaryasthan.org",
		ReplyTo:     "noreply@kaaryasthan.org",
		Subject:     "Test mail 1",
		HTMLContent: "This is a test mail from kaaryasthan.org",
	}
	err := ds.SendMail(eml)
	if err != nil {
		t.Error("Mail sendinging failed.", err)
	}
}
