package model

const (
	//Error 错误
	Error = "Error"
	//Broadcast 广播
	Broadcast = "Broadcast"
	//Join 新增
	Join = "Join"
)

//Message 消息
type Message struct {
	Path     string      `json:"path"`     //路径
	Kind     string      `json:"kind"`     //消息类型
	DataName string      `json:"dataName"` //数据名称
	Content  interface{} `json:"content"`  //消息类容
	MemberInfo
}

//NewErrorMessage 创建一个异常消息
func NewErrorMessage(path string, dataName string, content interface{}, m MemberInfo) Message {
	return Message{
		Path:       path,
		Kind:       Error,
		DataName:   dataName,
		Content:    content,
		MemberInfo: m,
	}
}
