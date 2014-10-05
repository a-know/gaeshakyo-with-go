package controller

import (
	"appengine"
	"encoding/json"
	"net/http"

	"model"
)

func WriteToGuestbook(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	message := r.FormValue("message")
	model.Save(c, message)
}

func GetMessageList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	messages := model.DescList(c)

	js, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
