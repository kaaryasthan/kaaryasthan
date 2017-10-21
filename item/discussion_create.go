package item

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func createDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	obj := new(Discussion)
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

	itm := Item{ID: obj.ItemID}
	if err := itm.Valid(); err != nil {
		log.Println("Couldn't validate item: ", err)
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

// Create creates new discussions
func (obj *Discussion) Create(usr user.User) error {
	err := db.DB.QueryRow(`INSERT INTO "discussions" (body, created_by, item_id) VALUES ($1, $2, $3) RETURNING id`,
		obj.Body, usr.ID, obj.ItemID).Scan(&obj.ID)
	return err
}
