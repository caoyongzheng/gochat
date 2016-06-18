package group

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
)

//Message 消息
type Message struct {
	Kind    string      `json:"kind"`
	Content interface{} `json:"content"`
}

//MemberInfo 成员信息
type MemberInfo struct {
	ID string `json:"id"`
}

//Member 成员
type Member struct {
	MemberInfo
	Group      *Group
	Connection *websocket.Conn
	Send       chan Message
}

//SendToClient 发送消息到客户端
func (m *Member) SendToClient() {
	for message := range m.Send {
		err := websocket.JSON.Send(m.Connection, message)
		if err != nil {
			break
		}
	}
}

//ReceiveFromClient 接受来自客户端的消息
func (m *Member) ReceiveFromClient() {
	for {
		var content string
		err := websocket.Message.Receive(m.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			log.Print(err)
			return
		}
		var message Message
		err = json.Unmarshal([]byte(content), &message)
		if err != nil {
			return
		}
		m.Group.Broadcast <- message
	}
}

//CloseChans 关闭所有管道
func (m *Member) CloseChans() {
	m.Connection.Close()
	close(m.Send)
}

//Group 群
type Group struct {
	ID          string                 `json:"id"`
	OrganizerID string                 `json:"organizerID"`
	Members     map[string]*Member     `json:"members"`
	Data        map[string]interface{} `json:"data"`
	CloseSign   chan bool
	Broadcast   chan Message
}

//Start 开启广播
func (g *Group) Start() {
	if _, isClose := <-g.Broadcast; isClose {
		g.Broadcast = make(chan Message)
	}
	for {
		select {
		case b := <-g.Broadcast:
			for _, member := range g.Members {
				member.Send <- b
			}
		case c := <-g.CloseSign:
			if c == true {
				close(g.Broadcast)
				close(g.CloseSign)
				return
			}
		}
	}
}

//Stop 停止广播
func (g *Group) Stop() {
	g.CloseSign <- true
}

//DeleteMember 删除用户
func (g *Group) DeleteMember(memberID string) {
	if member, ok := g.Members[memberID]; ok {
		member.CloseChans()
		delete(g.Members, memberID)
		log.Printf("Deleted %s", memberID)
		g.BroadcastMemberInfos()
	}
}

//BroadcastMemberInfos 广播成员信息
func (g *Group) BroadcastMemberInfos() {
	var memberInfos []MemberInfo
	for _, v := range g.Members {
		memberInfos = append(memberInfos, v.MemberInfo)
	}
	b, _ := json.Marshal(memberInfos)

	g.Broadcast <- Message{Kind: "Members", Content: string(b)}
}
