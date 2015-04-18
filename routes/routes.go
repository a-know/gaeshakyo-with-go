package routes

import (
	"net/http"

	"controller/auth"
	"controller/channel"
	"controller/cron"
	"controller/guestbook"
	"controller/memo"
	"controller/minutes"
	"controller/tq"
)

func init() {
	// chapter 2
	http.HandleFunc("/postGuestbook", guestbook.WriteToGuestbook)
	http.HandleFunc("/getGuestbook", guestbook.GetMessageList)

	// chapter 3.2
	http.HandleFunc("/postMinutes", minutes.Post)
	http.HandleFunc("/showMinutes", minutes.Show)
	http.HandleFunc("/postMemo", memo.Post)
	http.HandleFunc("/showMemo", memo.Show)

	// chapter 3.3
	http.HandleFunc("/auth", auth.Auth)

	// chapter 3.5
	http.HandleFunc("/tq/IncrementMemoCount", tq.IncrementMemoCount)

	// chapter 3.6
	http.HandleFunc("/cron/UpdateMemoCount", cron.UpdateMemoCount)

	// chapter 3.8
	http.HandleFunc("/channel", channel.CreateToken)

	// chapter 3.9
	http.HandleFunc("/deleteMinutes", minutes.Delete)

}
