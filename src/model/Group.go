package model

//Group 群
type Group struct {
	Path        string //路径
	Connections map[string]*Connection
}

//NewGroup 创建一个Group实例
func NewGroup(path string) *Group {
	return &Group{
		Path:        path,
		Connections: make(map[string]*Connection),
	}
}

//RemoveConnection 去除某一个连接
func (g *Group) RemoveConnection(connectionID string) {
	if _, ok := g.Connections[connectionID]; ok {
		delete(g.Connections, connectionID)
		g.broadcastMembers()
	}
}

//AddConnection 添加连接
func (g *Group) AddConnection(c *Connection) {
	g.Connections[c.ID] = c
	g.broadcastMembers()
}

//Broadcast 广播消息
func (g *Group) Broadcast(message Message) {
	dataName := message.DataName
	path := g.Path
	for _, c := range g.Connections {
		if c.IsListen(path, dataName) {
			c.Send <- message
		}
	}
}

//BroadcastMembers 广播人员
func (g *Group) broadcastMembers() {
	message := Message{
		Path:     g.Path,
		Kind:     Broadcast,
		DataName: "Members",
		Content:  g.GetMembers(),
	}
	g.Broadcast(message)
}

//GetMembers 获取所有连接
func (g *Group) GetMembers() []string {
	members := []string{}
	for _, c := range g.Connections {
		members = append(members, c.ID)
	}
	return members
}
