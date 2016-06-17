package main

import (
	"code.google.com/p/go.net/websocket"
	"handler"
	"net/http"
	"wscon"
)

func main() {
	//静态文件请求处理函数
	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/join", handler.Join)
	http.Handle("/chat", websocket.Handler(wscon.BuildConnection))
	wscon.InitActiveGroup()
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/"+r.URL.Path[len("/"):])
}
