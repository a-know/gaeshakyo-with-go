package routes

import (
	"net/http"

	"controller"
)

func init() {
	http.HandleFunc("/postGuestbook", controller.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", controller.GetMessageList)
}
