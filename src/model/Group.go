package model

//Group 群
type Group struct {
	Path        string //路径
	Connections map[string]*Connection
	Data        map[string]interface{}
}

//NewGroup 创建一个Group实例
func NewGroup(path string) *Group {
	return &Group{
		Path:        path,
		Connections: make(map[string]*Connection),
		Data:        make(map[string]interface{}),
	}
}

//RemoveConnection 去除某一个连接
func (g *Group) RemoveConnection(connectionID string) {
	if _, ok := g.Connections[connectionID]; ok {
		delete(g.Connections, connectionID)
		g.Data["Members"] = g.GetMembers()
	}
}

//AddConnection 添加连接
func (g *Group) AddConnection(c *Connection) {
	g.Connections[c.ID] = c
	g.Data["Members"] = g.GetMembers()
}

//Broadcast 广播消息
func (g *Group) Broadcast(message Message) {
	message.Status = Success
	for _, c := range g.Connections {
		c.Send <- message
	}
}

//BroadcastMembers 广播人员
func (g *Group) BroadcastMembers() {
	message := Message{
		Path:     g.Path,
		Kind:     Broadcast,
		DataName: "Members",
		Content:  g.Data["Members"],
	}
	g.Broadcast(message)
}

//BroadcastData 广播存储数据
func (g *Group) BroadcastData(dataName string) {
	if d, ok := g.Data[dataName]; ok {
		g.Broadcast(NewBroadcastMessage(g.Path, dataName, d, Member{}))
	}
}

//GetMembers 获取所有连接
func (g *Group) GetMembers() []string {
	members := []string{}
	for _, c := range g.Connections {
		members = append(members, c.ID)
	}
	return members
}
