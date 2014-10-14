package auth

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"net/http"

	"dto"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var url string
	var url_err error
	var dto dto.AuthDTO

	if u != nil {
		url, url_err = user.LogoutURL(c, "/statics/index.html")
		dto.LoggedIn  = true
		dto.LogoutURL = url
	} else {
		url, url_err = user.LoginURL(c, "/statics/index.html")
		dto.LoggedIn  = false
		dto.LogoutURL = url
	}
	if url_err != nil {
		http.Error(w, url_err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}