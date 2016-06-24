package model

const (
	//Error 错误
	Error = iota
	//Success 成功
	Success
)

const (
	//Broadcast 广播
	Broadcast = "Broadcast"
	//Unicast 单播
	Unicast = "Unicast"
	//Subscription 订阅
	Subscription = "Subscription"
	//ExitGroup 退出群
	ExitGroup = "ExitGroup"
)

//Message 消息
type Message struct {
	Path     string      `json:"path"`     //路径
	Kind     string      `json:"kind"`     //消息类型
	DataName string      `json:"dataName"` //数据名称
	Content  interface{} `json:"content"`  //消息类容
	Status   int         `json:"status"`   //消息状态
	Member
}

//NewErrorMessage 创建一个异常消息
func NewErrorMessage(path string, kind string, dataName string, content interface{}) Message {
	return Message{
		Path:     path,
		DataName: dataName,
		Content:  content,
		Status:   Error,
	}
}

//NewUnicastMessage 创建一个类型单播消息
func NewUnicastMessage(path string, dataName string, content interface{}, m Member) Message {
	return Message{
		Path:     path,
		Kind:     Unicast,
		DataName: dataName,
		Content:  content,
		Member:   m,
	}
}

//NewBroadcastMessage 创建一个广播类型消息
func NewBroadcastMessage(path string, dataName string, content interface{}, m Member) Message {
	return Message{
		Path:     path,
		Kind:     Broadcast,
		DataName: dataName,
		Content:  content,
		Member:   m,
	}
}

//NewSuccessMessage 创建一个成功消息
func NewSuccessMessage(path string, dataName string, content interface{}) Message {
	return Message{
		Path:     path,
		DataName: dataName,
		Content:  content,
		Status:   Success,
	}
}
