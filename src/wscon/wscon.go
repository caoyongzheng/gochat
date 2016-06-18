package wscon

import (
	"code.google.com/p/go.net/websocket"
	"global"
	"model"
	"service"
)

//BuildConnection 建立连接
func BuildConnection(ws *websocket.Conn) {
	id := ws.Request().URL.Query().Get("id")
	if id == "" {
		ws.Close()
		return
	}

	if _, ok := global.ActiveMembers[id]; ok {
		websocket.JSON.Send(ws, model.Message{Kind: "error", Content: "Id is already exist"})
		ws.Close()
		return
	}

	member := model.NewMember(id, ws)

	global.AddMember(member)

	go service.PushMessageToClient(member)
	service.HandleMessageFromClient(member)
}
