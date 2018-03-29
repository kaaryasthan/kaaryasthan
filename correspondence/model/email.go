package correspondence

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kaaryasthan/kaaryasthan/config"
)

// Email uses https://sendinblue.com to send mail
type Email struct {
	ToName      string
	ToEmail     string
	Sender      string
	ReplyTo     string
	Subject     string
	HTMLContent string
}

// SendMail sends mail
func (ds *Datastore) SendMail(eml *Email) error {
	url := "https://api.sendinblue.com/v3/smtp/email"

	plfmt := `{"sender":{"email":"%s"},"replyTo":{"email":"%s"},"to":[{"email":"%s","name":"%s"}],"htmlContent":"%s","subject":"%s"}`

	payload := strings.NewReader(fmt.Sprintf(plfmt, eml.Sender, eml.ReplyTo, eml.ToEmail, eml.ToName, eml.HTMLContent, eml.Subject))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", config.Config.SendinblueKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	c, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		return errors.New(res.Status + " " + string(c))
	}
	return nil
}
