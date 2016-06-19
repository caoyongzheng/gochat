package model

import (
	"code.google.com/p/go.net/websocket"
)

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
