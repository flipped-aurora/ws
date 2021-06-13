package biz

// IMessage 定义消息接口
type IMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	GetType() int32
	GetTo() string
}
