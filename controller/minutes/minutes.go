package minutes

import (
	"appengine"
	"appengine/datastore"
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

func Delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// ログインチェック
	if u == nil {
		http.Error(w, "you must to login", http.StatusUnauthorized)
		return
	}

	// 削除対象の議事録のキーを取得
	minutesKeyString := r.FormValue("delete")
	minutesKey, err := datastore.DecodeKey(minutesKeyString)
	if err != nil {
		http.Error(w, "not specified delete target", http.StatusBadRequest)
		return
	}

	// 議事録の取得
	var m minutes.Minutes
	err = datastore.Get(c, minutesKey, &m)
	if err != nil {
		c.Errorf("Error occurs: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m.Author.Email != u.Email {
		http.Error(w, "You are not author of this minutes", http.StatusForbidden)
		return
	}

	// tsv ファイルの作成とエンティティの削除
	fileName, export_err := minutes.ExportAsTsv(r, c, m)
	if export_err != nil {
		c.Errorf("Error occurs: %v", export_err)
		http.Error(w, export_err.Error(), http.StatusInternalServerError)
		return
	}
	minutes.Delete(c, m.Key)

	// ダウンロードURL をメールで送信する
	url, url_err := minutes.GetTsvUrl(r, c, fileName)
	if url_err != nil {
		c.Errorf("Error occurs: %v", url_err)
		http.Error(w, url_err.Error(), http.StatusInternalServerError)
		return
	}

	msg := &mail.Message{
		Sender:  "minutes@gaeshakyo-with-go.appspotmail.com",
		To:      []string{u.Email},
		Subject: "議事録「" + m.Title + "」の削除に伴い、内容が TSV ファイルに変換されました。",
		Body:    url,
	}
	if mail_err := mail.Send(c, msg); mail_err != nil {
		c.Errorf("Error occurs: %v", mail_err)
		http.Error(w, mail_err.Error(), http.StatusInternalServerError)
		return
	}
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
