package model

import (
	"code.google.com/p/go.net/websocket"
)

//Connection 成员
type Connection struct {
	ID string
	Member
	Conn    *websocket.Conn
	Send    chan Message
	Listens map[string]bool //key值为groupPath,
}

//Member 成员信息
type Member struct{}

//NewConnection 创建新成员
func NewConnection(id string, ws *websocket.Conn) *Connection {
	return &Connection{
		ID:      id,
		Conn:    ws,
		Send:    make(chan Message, 256),
		Listens: make(map[string]bool),
	}
}

//RemoveGroup 去除群
func (c *Connection) RemoveGroup(groupPath string) {
	if _, ok := c.Listens[groupPath]; ok {
		delete(c.Listens, groupPath)
	}
}

//AddListen 添加监听 return (bool)
func (c *Connection) AddListen(groupPath string) {
	c.Listens[groupPath] = true
}

//CloseChans 关闭所有通道
func (c *Connection) CloseChans() {
	c.Conn.Close()
	close(c.Send)
}
