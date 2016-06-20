package model

//GroupTree 群组树
type GroupTree struct {
	Name       string                //名字
	Path       string                //路径
	GroupTrees map[string]*GroupTree //子树群
	Members    map[string]*Member    //成员
}

//NewGroupTree 创建一个GroupTree实例
func NewGroupTree(name string, path string) *GroupTree {
	return &GroupTree{
		Name:       name,
		Path:       path,
		GroupTrees: make(map[string]*GroupTree),
		Members:    make(map[string]*Member),
	}
}

//RemoveMember 删除成员
func (groupTree *GroupTree) RemoveMember(memberID string) {
	if _, ok := groupTree.Members[memberID]; ok {
		delete(groupTree.Members, memberID)
	}
}
