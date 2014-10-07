package routes

import (
	"net/http"

	"controller/guestbook"
)

func init() {
	http.HandleFunc("/postGuestbook", guestbook.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", guestbook.GetMessageList)
}
