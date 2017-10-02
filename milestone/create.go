package milestone

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]jsonapi.Data)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}
	s := New(payload["data"])
	id, err := s.create()
	if err != nil {
		log.Println("Unable save data: ", err)
		return
	}
	tmpData := payload["data"]
	tmpData.ID = strconv.Itoa(id)
	payload["data"] = tmpData
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Unable marshal data: ", err)
		return
	}
	w.Write(b)
}

func (obj *Schema) create() (int, error) {
	var id int
	err := db.DB.QueryRow(`INSERT INTO "milestones" (name, description) VALUES ($1, $2) RETURNING id`, obj.Name, obj.Description).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
