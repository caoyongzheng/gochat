package global

import (
	"fmt"
	"model"
	"strings"
	"sync"
)

//Mu 锁
var Mu sync.RWMutex

//ActiveGroupTree 当前活动群
var ActiveGroupTree *model.GroupTree

//System 系统
var System model.MemberInfo

func init() {
	ActiveGroupTree = model.NewGroupTree("system", "system")
	System = model.MemberInfo{
		ID: "SYSTEM",
	}
}

//GetAndAddGroupTree 获取groupTree
func GetAndAddGroupTree(path string) (*model.GroupTree, error) {
	pathItems := strings.Split(path, ".")
	if pathItems[0] != "system" {
		return nil, fmt.Errorf("root paths is wrong,expected system,but is %s", pathItems[0])
	}
	var group *model.GroupTree
	group = ActiveGroupTree
	for k, item := range pathItems {
		if k == 0 {
			continue
		}
		if group.GroupTrees == nil {
			group.GroupTrees = make(map[string]*model.GroupTree)
		}
		if g, ok := group.GroupTrees[item]; ok {
			group = g
		} else {
			group.GroupTrees[item] = model.NewGroupTree(item, path)
			group = group.GroupTrees[item]
		}
	}
	return group, nil
}

//AddMember 添加新成员
func AddMember(member *model.Member) {
	defer Mu.Unlock()
	Mu.Lock()
	ActiveGroupTree.Members[member.ID] = member
}
