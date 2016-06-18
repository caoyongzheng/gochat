package model

import (
	"encoding/json"
	"log"
)

//Group 群
type Group struct {
	ID          string                 `json:"id"`
	OrganizerID string                 `json:"organizerID"`
	Members     map[string]*Member     `json:"members"`
	Data        map[string]interface{} `json:"data"`
	CloseSign   chan bool
	Broadcast   chan Message
	IsClose     bool
}

//Start 开启广播
func (g *Group) Start() {
	if g.IsClose {
		g.Broadcast = make(chan Message)
		g.CloseSign = make(chan bool)
		g.IsClose = false
	}
	for {
		select {
		case b := <-g.Broadcast:
			for _, member := range g.Members {
				member.Send <- b
			}
		case c := <-g.CloseSign:
			if c == true {
				close(g.Broadcast)
				close(g.CloseSign)
				return
			}
		}
	}
}

//Stop 停止广播
func (g *Group) Stop() {
	g.IsClose = true
	g.CloseSign <- true
}

//DeleteMember 删除用户
func (g *Group) DeleteMember(memberID string) {
	if _, ok := g.Members[memberID]; ok {
		delete(g.Members, memberID)
		log.Printf("Deleted %s", memberID)
		g.BroadcastMemberInfos()
	}
}

//BroadcastMemberInfos 广播成员信息
func (g *Group) BroadcastMemberInfos() {
	var memberInfos []MemberInfo
	for _, v := range g.Members {
		memberInfos = append(memberInfos, v.MemberInfo)
	}
	b, _ := json.Marshal(memberInfos)

	g.Broadcast <- Message{Kind: "Members", Content: string(b)}
}
