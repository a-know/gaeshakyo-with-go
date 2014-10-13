package routes

import (
	"net/http"

	"controller/guestbook"
	"controller/minutes"
)

func init() {
	// chapter 2
	http.HandleFunc("/postGuestbook", guestbook.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", guestbook.GetMessageList)

	// chapter 3
	http.HandleFunc("/postMinutes", minutes.Post)
	http.HandleFunc("/showMinutes", minutes.Show)
}
