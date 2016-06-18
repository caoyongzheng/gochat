package service

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"global"
	"log"
	"model"
)

//PushMessageToClient 推送消息到客户端
func PushMessageToClient(member *model.Member) {
	for message := range member.Send {
		err := websocket.JSON.Send(member.Connection, message)
		if err != nil {
			break
		}
	}
}

//HandleMessageFromClient 处理来自客户端的消息
func HandleMessageFromClient(member *model.Member) {
	for {
		var content string
		err := websocket.Message.Receive(member.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			log.Print(err)
			HandleMemberOffLine(member)
			return
		}
		var message model.Message
		err = json.Unmarshal([]byte(content), &message)
		if err != nil {
			return
		}
		if group, ok := member.Groups[message.GroupID]; ok {
			if !group.IsClose {
				group.Broadcast <- message
			} else {
				member.Send <- model.Message{Kind: "error", Content: message.GroupID + "is closed"}
			}
		} else {
			member.Send <- model.Message{Kind: "error", Content: message.GroupID + "is not exist"}
		}
	}
}

//HandleMemberOffLine 处理用户下线
func HandleMemberOffLine(member *model.Member) {
	member.OffLine()
	global.DeleteMember(member.ID)
}
