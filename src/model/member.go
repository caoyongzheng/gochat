package model

import (
	"code.google.com/p/go.net/websocket"
)

//Message 消息
type Message struct {
	Kind    string      `json:"kind"`
	GroupID string      `json:"groupID"`
	Content interface{} `json:"content"`
}

//MemberInfo 成员信息
type MemberInfo struct {
	ID string `json:"id"`
}

//Member 成员
type Member struct {
	MemberInfo
	Groups     map[string]*Group
	Connection *websocket.Conn
	Send       chan Message
}

//NewMember 创建一个成员实例
func NewMember(id string, ws *websocket.Conn) *Member {
	return &Member{
		MemberInfo: MemberInfo{
			ID: id,
		},
		Connection: ws,
		Send:       make(chan Message, 256),
	}
}

//OffLine 用户下线
func (m *Member) OffLine() {
	for _, group := range m.Groups {
		m.ExitGroup(group.ID)
	}
	m.Connection.Close()
	close(m.Send)
}

//ExitGroup 用户退出群
func (m *Member) ExitGroup(groupID string) {
	if group, ok := m.Groups[groupID]; ok {
		group.DeleteMember(m.ID)
		delete(m.Groups, groupID)
	}
}
