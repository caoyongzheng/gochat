package service

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
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
		if err != nil {
			HandleOffLine(connection)
			return
		}
		var message model.Message
		err = json.Unmarshal([]byte(content), &message)
		if err != nil {
			HandleOffLine(connection)
			return
		}
		message.Member = connection.Member
		message.Status = model.Success
		handleMessageFromClient(message, connection)
	}
}

func handleMessageFromClient(message model.Message, connection *model.Connection) {
	switch message.Kind {
	case model.Subscription:
		handleSubscription(message, connection)
	case model.ExitGroup:
		handleExitGroup(message, connection)
	case model.Broadcast:
		handleBroadcast(message, connection)
	default:
		message.Status = model.Error
		message.Content = "Kind:" + message.Kind + "is unrecognized"
		connection.Send <- message
	}
}

func handleSubscription(message model.Message, connection *model.Connection) {
	path := message.Path
	dataName := message.DataName
	group, ok := global.ActiveGroups[path]

	if !ok {
		group = model.NewGroup(path)
		global.AddGroup(group)
	}

	connection.AddListen(group.Path)
	group.AddConnection(connection)

	message.Content = fmt.Sprintf("Success to listen %s:%s", path, dataName)
	connection.Send <- message

	group.BroadcastData("Members")
}

func handleExitGroup(message model.Message, connection *model.Connection) {
	path := message.Path
	if group, ok := global.ActiveGroups[path]; ok {
		connection.RemoveGroup(path)
		group.RemoveConnection(connection.ID)
	}

	message.Status = model.Success
	message.Content = fmt.Sprintf("Success to exit group %s", path)
	connection.Send <- message
}

func handleBroadcast(message model.Message, connection *model.Connection) {
	path := message.Path
	group, ok := global.ActiveGroups[path]

	if !ok {
		message.Status = model.Error
		message.Content = "path:" + path + "is not exist"
		connection.Send <- message
		return
	}

	group.Broadcast(message)
}

//HandleOffLine 处理用户下线
func HandleOffLine(connection *model.Connection) {
	defer global.Mu.Unlock()
	global.Mu.Lock()

	id := connection.ID
	for k := range connection.Listens {
		if g, ok := global.ActiveGroups[k]; ok {
			g.RemoveConnection(id)
			g.BroadcastData("Members")
		}
	}

	if _, ok := global.ActiveConnections[id]; ok {
		delete(global.ActiveConnections, id)
	}

	connection.CloseChans()
	log.Printf("Connection %s is offline", id)
}
