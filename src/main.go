package main

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"wscon"
)

func main() {
	http.Handle("/chat", websocket.Handler(wscon.BuildConnection))
	wscon.InitActiveGroup()
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
