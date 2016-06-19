package wscon

import (
	"code.google.com/p/go.net/websocket"
	"global"
	"log"
	"model"
	"service"
)

//BuildConnection 建立连接
func BuildConnection(ws *websocket.Conn) {
	id := ws.Request().URL.Query().Get("id")
	if id == "" {
		websocket.JSON.Send(ws, model.NewErrorMessage("system", "Error", "id cannot be null", global.System))
		ws.Close()
		return
	}

	if _, ok := global.ActiveGroupTree.Members[id]; ok {
		websocket.JSON.Send(ws, model.NewErrorMessage("system", "Error", "Id:"+id+" is already exist", global.System))
		ws.Close()
		return
	}

	member := model.NewMember(id, ws, global.ActiveGroupTree)
	global.AddMember(member)
	log.Printf("Member %s is online", id)
	go service.SendMessageToClient(member)
	service.ReceiveMessageFromClient(member)
}
