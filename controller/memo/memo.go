package memo

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"net/http"

	"model/memo"
)

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	memoString := r.FormValue("memo")
	minutesKeyString := r.FormValue("minutes")
	minutesKey, err := datastore.DecodeKey(minutesKeyString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if memoString == "" {
		http.Error(w, "memo is not set", http.StatusBadRequest)
		return
	}

	_, err = memo.SaveAs(c, minutesKey ,memoString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	minutesKeyString := r.FormValue("minutes")
	minutesKey, decode_err := datastore.DecodeKey(minutesKeyString)

	if decode_err != nil {
		http.Error(w, decode_err.Error(), http.StatusBadRequest)
		return
	}

	memo_list, list_err := memo.AscList(c, minutesKey)

	js, list_err := json.Marshal(memo_list)
	if list_err != nil {
		http.Error(w, list_err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
