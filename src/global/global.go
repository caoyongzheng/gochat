package global

import (
	"model"
	"sync"
)

var mu sync.RWMutex

//ActiveGroups 当前活动群
var ActiveGroups map[string]*model.Group

//ActiveMembers 当前活跃用户
var ActiveMembers map[string]*model.Member

//AddMember 添加成员
func AddMember(member *model.Member) {
	defer mu.Unlock()
	mu.Lock()
	ActiveMembers[member.ID] = member
}

//DeleteMember 删除用户
func DeleteMember(memberID string) {
	if _, ok := ActiveMembers[memberID]; ok {
		delete(ActiveMembers, memberID)
	}
}
