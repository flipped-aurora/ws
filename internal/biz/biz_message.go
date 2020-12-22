package biz

type Message struct {
	Type int32  `json:"type"`
	Time int64  `json:"time"`
	From string `json:"From"`
	To   string `json:"to" binding:"required"`
	Data []byte `json:"data" binding:"required"`
}

// IMessage: 定义消息接口
type IMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
