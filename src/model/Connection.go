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
	Listens map[string][]string //key值为groupPath,value监听的数据
}

//Member 成员信息
type Member struct{}

//NewConnection 创建新成员
func NewConnection(id string, ws *websocket.Conn) *Connection {
	return &Connection{
		ID:      id,
		Conn:    ws,
		Send:    make(chan Message, 256),
		Listens: make(map[string][]string),
	}
}

//AddListen 添加监听 return (bool)
func (c *Connection) AddListen(groupPath string, dataName string) bool {
	if ds, ok := c.Listens[groupPath]; ok {
		for _, v := range ds {
			if v == dataName {
				return false
			}
		}
		c.Listens[groupPath] = append(ds, dataName)
		return false
	}
	c.Listens[groupPath] = []string{dataName}
	return true
}

//CancelListen 取消监听
func (c *Connection) CancelListen(groupPath string, dataName string) bool {
	if ds, ok := c.Listens[groupPath]; ok {
		for k, v := range ds {
			if v == dataName {
				ds[k] = ds[len(ds)-1]
				ds = ds[:len(ds)-1]
				break
			}
		}
		if len(ds) == 0 {
			delete(c.Listens, groupPath)
			return true
		}
		return false
	}
	return false
}

//IsListen 是否监听
func (c *Connection) IsListen(groupPath, dataName string) bool {
	if ds, ok := c.Listens[groupPath]; ok {
		for _, v := range ds {
			if v == dataName {
				return true
			}
		}
	}
	return false
}

//CloseChans 关闭所有通道
func (c *Connection) CloseChans() {
	c.Conn.Close()
	close(c.Send)
}
