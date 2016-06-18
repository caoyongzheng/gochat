package wscon

import (
	"code.google.com/p/go.net/websocket"
	"model/group"
)

//ActiveGroup 当前活动群
var ActiveGroup *group.Group

//BuildConnection 建立连接
func BuildConnection(ws *websocket.Conn) {
	email := ws.Request().URL.Query().Get("id")
	if email == "" {
		ws.Close()
		return
	}
	if _, ok := ActiveGroup.Members[email]; ok {
		ws.Close()
		return
	}

	member := &group.Member{
		MemberInfo: group.MemberInfo{
			ID: email,
		},
		Group:      ActiveGroup,
		Connection: ws,
		Send:       make(chan group.Message, 256),
	}

	ActiveGroup.Members[email] = member
	ActiveGroup.BroadcastMemberInfos()

	go member.SendToClient()
	member.ReceiveFromClient()
	ActiveGroup.DeleteMember(email)
}

//InitActiveGroup 初始化房间
func InitActiveGroup() {
	ActiveGroup = &group.Group{
		Members:   make(map[string]*group.Member),
		Broadcast: make(chan group.Message),
		CloseSign: make(chan bool),
	}
	go ActiveGroup.Start()
}
