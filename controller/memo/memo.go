package memo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"net/http"

	"model/memo"
)

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		http.Error(w, "you must to login", http.StatusUnauthorized)
		return
	}

	memoString := r.FormValue("memo")

	if memoString == "" {
		http.Error(w, "memo is not set", http.StatusBadRequest)
		return
	}

	minutesKeyString := r.FormValue("minutes")
	minutesKey, err := datastore.DecodeKey(minutesKeyString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memoKey, put_err := memo.SaveAs(c, minutesKey, memoString, u)
	if put_err != nil {
		http.Error(w, put_err.Error(), http.StatusInternalServerError)
		return
	}
	_ = memo.PushNotification(c, memoKey) // notofication に失敗しても特にエラーとはしない
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
