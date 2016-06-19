package service

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"global"
	"log"
	"model"
)

//SendMessageToClient 推送消息到客户端
func SendMessageToClient(member *model.Member) {
	for message := range member.Send {
		err := websocket.JSON.Send(member.Connection, message)
		if err != nil {
			break
		}
	}
}

//ReceiveMessageFromClient 处理来自客户端的消息
func ReceiveMessageFromClient(member *model.Member) {
	for {
		var content string
		err := websocket.Message.Receive(member.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			HandleMemberOffLine(member)
			return
		}
		var message model.Message
		err = json.Unmarshal([]byte(content), &message)
		if err != nil {
			HandleMemberOffLine(member)
			return
		}
		message.MemberInfo = member.MemberInfo
		handleMessageFromClient(message, member)
	}
}

func handleMessageFromClient(message model.Message, member *model.Member) {
	path := message.Path
	switch message.Kind {
	case model.Join:
		group, err := global.GetAndAddGroupTree(path)
		if err != nil {
			member.Send <- model.NewErrorMessage(path, message.DataName, err, global.System)
		}
		group.Members[member.ID] = member
		member.GroupTrees[group.Path] = group
	case model.Broadcast:
		if groupTree, ok := member.GroupTrees[path]; ok {
			for _, m := range groupTree.Members {
				m.Send <- message
			}
		} else {
			member.Send <- model.NewErrorMessage(path, message.DataName, "path:"+path+"is not exist", global.System)
		}
	default:
		member.Send <- model.NewErrorMessage(path, message.DataName, "Kind:"+message.Kind+"is unrecognized", global.System)
	}
}

//HandleMemberOffLine 处理用户下线
func HandleMemberOffLine(member *model.Member) {
	defer global.Mu.Unlock()
	global.Mu.Lock()
	member.OffLine()
	log.Printf("Member %s is offline", member.ID)
}
