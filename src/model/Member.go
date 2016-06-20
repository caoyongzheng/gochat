package model

import (
	"code.google.com/p/go.net/websocket"
)

//Connection 成员
type Connection struct {
	ID string
	Member
	Conn       *websocket.Conn
	Send       chan Message
	GroupTrees map[string]*GroupTree
}

//Member 成员信息
type Member struct{}

//NewConnection 创建新成员
func NewConnection(id string, ws *websocket.Conn, groupTree *GroupTree) *Connection {
	return &Connection{
		Member:     Member{ID: id},
		Conn:       ws,
		Send:       make(chan Message, 256),
		GroupTrees: map[string]*GroupTree{groupTree.Path: groupTree},
	}
}

//OffLine yong
func (c *Connection) OffLine() {
	for _, group := range c.GroupTrees {
		group.RemoveConnection(c.ID)
	}
	c.GroupTrees = nil
	c.Conn.Close()
	close(c.Send)
}
