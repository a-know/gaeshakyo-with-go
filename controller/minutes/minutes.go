package minutes

import (
	"appengine"
	"appengine/mail"
	"appengine/user"
	"encoding/json"
	"net/http"

	"model/minutes"
)

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		http.Error(w, "you must to login", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")

	if title != "" {
		minutesKey, err := minutes.SaveAs(c, title, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendMail(r, c, minutesKey.Encode())
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

func sendMail(r *http.Request, c appengine.Context, minutesKeyString string) {
	// メールの本文に入れる議事録リンクの文字列を作成する
	scheme := r.URL.Scheme
	host := r.URL.Host
	link := scheme + "://" + host + "/statics/minutes.html?minutes=" + minutesKeyString

	msg := &mail.Message{
		Sender:  "minutes@gaeshakyo-with-go.appspotmail.com",
		Subject: "新しい議事録が追加されました",
		Body:    link,
	}
	if err := mail.SendToAdmins(c, msg); err != nil {
		c.Errorf("the email failed to send: %v", err)
	}
}
