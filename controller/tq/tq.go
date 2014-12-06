package tq

import (
	"appengine"
	"appengine/datastore"
	"model/minutes"
	"net/http"
)

func IncrementMemoCount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	minutesKeyString := r.FormValue("minutesKey")
	minutesKey, decode_err := datastore.DecodeKey(minutesKeyString)
	if decode_err != nil {
		http.Error(w, "irregal key string", http.StatusBadRequest)
	}
	minutes.IncrementMemoCount(c, minutesKey)
}
