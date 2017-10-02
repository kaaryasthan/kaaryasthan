package discussion

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item, err := strconv.Atoi(vars["item"])
	if err != nil {
		log.Println("Item is not an int value: ", err)
		return
	}

	payload := make(map[string]Data)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		log.Println("Unable to decode body: ", err)
		return
	}
	s, err := New(payload["data"])
	if err != nil {
		log.Println("Unable to create object: ", err)
		return
	}
	if item != s.Item {
		log.Println("Item ID in URL and payload is not matching")
		return
	}
	id, err := s.Create()
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
