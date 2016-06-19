package main

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"wscon"
)

func main() {
	http.Handle("/", websocket.Handler(wscon.BuildConnection))
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
