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
		message.MemberInfo = connection.MemberInfo
		handleMessageFromClient(message, connection)
	}
}

func handleMessageFromClient(message model.Message, connection *model.Connection) {
	path := message.Path
	switch message.Kind {
	case model.Join:
		group, err := global.GetAndAddGroupTree(path)
		if err != nil {
			connection.Send <- model.NewErrorMessage(path, message.DataName, err, global.System)
		}
		group.Connections[connection.ID] = connection
		connection.GroupTrees[group.Path] = group
	case model.Broadcast:
		if groupTree, ok := connection.GroupTrees[path]; ok {
			for _, m := range groupTree.Connections {
				m.Send <- message
			}
		} else {
			connection.Send <- model.NewErrorMessage(path, message.DataName, "path:"+path+"is not exist", global.System)
		}
	default:
		connection.Send <- model.NewErrorMessage(path, message.DataName, "Kind:"+message.Kind+"is unrecognized", global.System)
	}
}

//HandleConnectionOffLine 处理用户下线
func HandleConnectionOffLine(connection *model.Connection) {
	defer global.Mu.Unlock()
	global.Mu.Lock()
	connection.OffLine()
	log.Printf("Connection %s is offline", connection.ID)
}
