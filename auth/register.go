package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/config"
	correspondence "github.com/kaaryasthan/kaaryasthan/correspondence/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// RegisterHandler register user
func (c *Controller) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	usr := new(user.User)
	if err := jsonapi.UnmarshalPayload(r.Body, usr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.uds.Create(usr)
	if err != nil {
		log.Println("Unable save data: ", err)
		return
	}
	subject := "Kaaryasthan Login Email Verification"
	toName := usr.Username
	toEmail := usr.Email
	sender := config.Config.EmailSender
	replyTo := config.Config.EmailReplyTo
	htmlContentFormat := `Hi,<br/><br/> Thank you for registering at Kaaryasthan! <br/><br/> To confirm your email address, please click this link:<br/> %s <br/><br/>Thanks,<br/> Kaaryasthan Team `
	randomkey := usr.EmailVerificationCode
	htmlContent := fmt.Sprintf(htmlContentFormat, config.Config.BaseURL+"/email?key="+randomkey)
	email := &correspondence.Email{
		ToName:      toName,
		ToEmail:     toEmail,
		Sender:      sender,
		ReplyTo:     replyTo,
		Subject:     subject,
		HTMLContent: htmlContent,
	}
	usr.Password = ""
	if err := c.cds.SendMail(email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, usr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
