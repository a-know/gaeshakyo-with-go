package channel

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"encoding/json"
	"net/http"
	"time"

	channel_model "model/channel"

	"code.google.com/p/go-uuid/uuid"
)

func CreateToken(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	minutesKeyString := r.FormValue("minutes")
	minutesKey, err := datastore.DecodeKey(minutesKeyString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client_id := uuid.New()

	tok, err := channel.Create(c, client_id)
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}

	mc_key := datastore.NewKey(c, "minutes_channel", client_id, 0, nil)
	mc := channel_model.MinutesChannel{
		Key:        mc_key,
		Token:      tok,
		MinutesKey: minutesKey,
		CreatedAt:  time.Now(),
	}

	// put
	_, put_err := datastore.Put(c, mc_key, &mc)

	if put_err != nil {
		http.Error(w, put_err.Error(), http.StatusInternalServerError)
		return
	}

	js, js_err := json.Marshal(mc)
	if js_err != nil {
		http.Error(w, js_err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
