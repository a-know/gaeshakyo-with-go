package minutes

import (
	"appengine"
	"encoding/json"
	"net/http"

	"model/minutes"
)

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	title := r.FormValue("title")

	if title != "" {
		_, err := minutes.SaveAs(c, title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "title is not set", http.StatusBadRequest)
		return
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	minutes_list, err := minutes.DescList(c)

	js, err := json.Marshal(minutes_list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
