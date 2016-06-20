package model

//GroupTree 群组树
type GroupTree struct {
	Name        string                 //名字
	Path        string                 //路径
	GroupTrees  map[string]*GroupTree  //子树群
	Connections map[string]*Connection //连接
}

//NewGroupTree 创建一个GroupTree实例
func NewGroupTree(name string, path string) *GroupTree {
	return &GroupTree{
		Name:        name,
		Path:        path,
		GroupTrees:  make(map[string]*GroupTree),
		Connections: make(map[string]*Connection),
	}
}

//RemoveConnection 删除成员
func (groupTree *GroupTree) RemoveConnection(connectionID string) {
	if _, ok := groupTree.Connections[connectionID]; ok {
		delete(groupTree.Connections, connectionID)
	}
}
