package routes

import (
	"net/http"

	"controller"
	"hello"
)

func init() {
	http.HandleFunc("/", hello.Handler)
	http.HandleFunc("/postGuestbook", controller.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", controller.GetMessageList)
}
