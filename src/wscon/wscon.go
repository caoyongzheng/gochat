package wscon

import (
	"code.google.com/p/go.net/websocket"
	"global"
	"log"
	"model"
	"myutils"
	"service"
)

//BuildConnection 建立连接
func BuildConnection(ws *websocket.Conn) {
	id, err := myutils.UniqueID()
	if err != nil {
		websocket.JSON.Send(ws, model.NewErrorMessage("system", "Error", "Id:"+id+" is already exist", global.System))
		ws.Close()
		return
	}
	connection := model.NewConnection(id, ws)

	global.AddConnection(connection)
	log.Printf("Connection %s is online", id)

	go service.SendMessageToClient(connection)
	service.ReceiveMessageFromClient(connection)
}
