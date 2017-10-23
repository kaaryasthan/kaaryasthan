package item

import (
	"log"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// Show an item
func (ds *Datastore) Show(itm *Item) error {
	err := db.DB.QueryRow(`SELECT id, title, description, open_state,
			project_id, lock_conversation, created_by, updated_by, assignees,
			subscribers, labels, created_at, updated_at FROM "items"
			WHERE num=$1 AND deleted_at IS NULL`,
		itm.Number).Scan(&itm.ID, &itm.Title, &itm.Description,
		&itm.OpenState, &itm.ProjectID, &itm.LockConversation, &itm.CreatedBy,
		&itm.UpdatedBy, &itm.Assignees, &itm.Subscribers, &itm.Labels,
		&itm.CreatedAt, &itm.UpdatedAt)
	return err
}

// ShowItemHandler shows item
func (c *Controller) ShowItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: "+usr.ID, err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	num := vars["number"]

	number, err := strconv.Atoi(num)
	if err != nil {
		log.Println("Invalid number: "+num, err)
		http.Error(w, "Invalid number: "+num, http.StatusUnauthorized)
		return
	}

	itm := &Item{Number: number}
	if err := c.ds.Show(itm); err != nil {
		log.Println("Couldn't find item: ", number, err)
		http.Error(w, "Couldn't find item: "+string(number), http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, itm); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
