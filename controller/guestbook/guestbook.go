package guestbook

import (
	"appengine"
	"encoding/json"
	"net/http"

	"model/guestbook"
)

func WriteToGuestbook(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	message := r.FormValue("message")

	err := guestbook.Save(c, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetMessageList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	messages, err := guestbook.DescList(c)

	js, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
