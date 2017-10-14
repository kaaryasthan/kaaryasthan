package item

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	obj := new(Comment)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := user.User{ID: userID}
	if err := usr.Valid(); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	disc := Discussion{ID: obj.DiscussionID}
	if err := disc.Valid(); err != nil {
		log.Println("Couldn't validate discussion: ", err)
		return
	}

	err := obj.Create(usr)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create creates new comments
func (obj *Comment) Create(usr user.User) error {
	err := db.DB.QueryRow(`INSERT INTO "comments" (body, created_by, discussion_id) VALUES ($1, $2, $3) RETURNING id`,
		obj.Body, usr.ID, obj.DiscussionID).Scan(&obj.ID)
	if err != nil {
		return err
	}
	return nil
}
