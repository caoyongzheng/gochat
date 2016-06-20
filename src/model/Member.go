package model

import (
	"code.google.com/p/go.net/websocket"
)

//Member 成员
type Member struct {
	MemberInfo
	Connection *websocket.Conn
	Send       chan Message
	GroupTrees map[string]*GroupTree
}

//MemberInfo 成员信息
type MemberInfo struct {
	ID string
}

//NewMember 创建新成员
func NewMember(id string, ws *websocket.Conn, groupTree *GroupTree) *Member {
	return &Member{
		MemberInfo: MemberInfo{ID: id},
		Connection: ws,
		Send:       make(chan Message, 256),
		GroupTrees: map[string]*GroupTree{groupTree.Path: groupTree},
	}
}

//OffLine yong
func (m *Member) OffLine() {
	for _, group := range m.GroupTrees {
		group.RemoveMember(m.ID)
	}
	m.GroupTrees = nil
	m.Connection.Close()
	close(m.Send)
}
