package routes

import (
	"net/http"

	"controller/guestbook"
	"controller/minutes"
	"controller/memo"
	"controller/auth"
)

func init() {
	// chapter 2
	http.HandleFunc("/postGuestbook", guestbook.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", guestbook.GetMessageList)

	// chapter 3.2
	http.HandleFunc("/postMinutes", minutes.Post)
	http.HandleFunc("/showMinutes", minutes.Show)
	http.HandleFunc("/postMemo",    memo.Post)
	http.HandleFunc("/showMemo",    memo.Show)

	// chapter 3.3
	http.HandleFunc("/auth", auth.Auth)

}
