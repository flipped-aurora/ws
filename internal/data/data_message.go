package data

import "encoding/json"

type Message struct {
	Type int32  `json:"type"`
	Time int64  `json:"time"`
	From string `json:"From"`
	To   string `json:"to" binding:"required"`
	Data []byte `json:"data" binding:"required"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// GetType 获取状态信息
func (m *Message) GetType() int32 {
	return m.Type
}

// GetTo 获取接收人
func (m *Message) GetTo() string {
	return m.To
}
