package global

import (
	"model"
	"sync"
)

//Mu 锁
var Mu sync.RWMutex

//ActiveConnections 当前连接
var ActiveConnections map[string]*model.Connection

//ActiveGroups 当前活动群
var ActiveGroups map[string]*model.Group

//System 系统
var System model.Member

func init() {
	ActiveConnections = make(map[string]*model.Connection)
	ActiveGroups = make(map[string]*model.Group)
}

//AddConnection 添加连接
func AddConnection(c *model.Connection) {
	defer Mu.Unlock()
	Mu.Lock()
	ActiveConnections[c.ID] = c
}

//AddGroup 添加群
func AddGroup(g *model.Group) {
	defer Mu.Unlock()
	Mu.Lock()
	ActiveGroups[g.Path] = g
}
