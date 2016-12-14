package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// NewTestServer helper for testing
func NewTestServer(path, method string, fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	hts := httptest.NewServer(func() *mux.Router {
		rt := mux.NewRouter()
		rt.HandleFunc(path, fn).Methods(method)
		return rt
	}())
	return hts
}
