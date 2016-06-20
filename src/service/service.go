package service

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"global"
	"log"
	"model"
)

//SendMessageToClient 推送消息到客户端
func SendMessageToClient(connection *model.Connection) {
	for message := range connection.Send {
		err := websocket.JSON.Send(connection.Conn, message)
		if err != nil {
			break
		}
	}
}

//ReceiveMessageFromClient 处理来自客户端的消息
func ReceiveMessageFromClient(connection *model.Connection) {
	for {
		var content string
		err := websocket.Message.Receive(connection.Conn, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			HandleConnectionOffLine(connection)
			return
		}
		var message model.Message
		err = json.Unmarshal([]byte(content), &message)
		if err != nil {
			HandleConnectionOffLine(connection)
			return
		}
		message.Member = connection.Member
		handleMessageFromClient(message, connection)
	}
}

func handleMessageFromClient(message model.Message, connection *model.Connection) {
	path := message.Path
	dataName := message.DataName
	switch message.Kind {
	case model.Listen:
		group, ok := global.ActiveGroups[path]
		if !ok {
			group = model.NewGroup(path)
			global.AddGroup(group)
		}
		if dataName == "" {
			connection.Send <- model.NewErrorMessage(path, dataName, "dataName should not be empty", global.System)
		}
		if connection.AddListen(group.Path, dataName) {
			group.AddConnection(connection)
		}
	case model.Broadcast:
		group, ok := global.ActiveGroups[path]
		if !ok {
			connection.Send <- model.NewErrorMessage(path, message.DataName, "path:"+path+"is not exist", global.System)
		}
		if dataName == "" {
			connection.Send <- model.NewErrorMessage(path, dataName, "dataName should not be empty", global.System)
		}
		group.Broadcast(message)
	default:
		connection.Send <- model.NewErrorMessage(path, message.DataName, "Kind:"+message.Kind+"is unrecognized", global.System)
	}
}

//HandleConnectionOffLine 处理用户下线
func HandleConnectionOffLine(connection *model.Connection) {
	defer global.Mu.Unlock()
	global.Mu.Lock()
	id := connection.ID
	for k := range connection.Listens {
		if g, ok := global.ActiveGroups[k]; ok {
			g.RemoveConnection(id)
		}
	}
	if _, ok := global.ActiveConnections[id]; ok {
		delete(global.ActiveConnections, id)
	}
	connection.CloseChans()
	log.Printf("Connection %s is offline", id)
}
