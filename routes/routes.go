package routes

import (
	"net/http"

	"hello"
)

func init() {
	http.HandleFunc("/", hello.Handler)
}
