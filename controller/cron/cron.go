package cron

import (
	"appengine"
	"model/minutes"
	"net/http"
)

func UpdateMemoCount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	minutesKeyList, err := minutes.QueryForUpdateMemoCount(c)
	if err != nil {
		http.Error(w, "can't get targets for update memo count", http.StatusInternalServerError)
	}

	for _, minutesKey := range minutesKeyList {
		err = minutes.UpdateMemoCount(c, minutesKey)
		if err != nil {
			http.Error(w, "can't update target for update memo count", http.StatusInternalServerError)
		}
	}
}
