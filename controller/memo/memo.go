package memo

import (
	"appengine"
	"appengine/datastore"
	// "encoding/json"
	"net/http"

	"model/memo"
)

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	memoString := r.FormValue("memo")
	minutesKeyString := r.FormValue("minutes")
	minutesKey, _ := datastore.DecodeKey(minutesKeyString)

	if memoString != "" {
		_, err := memo.SaveAs(c, minutesKey ,memoString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "memo is not set", http.StatusBadRequest)
		return
	}
}