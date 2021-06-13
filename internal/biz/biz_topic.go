package biz

// ITopic 订阅表
type ITopic interface {
	// CreateTopic 创建topic
	CreateTopic(string)
	// DeleteTopic 删除topic
	DeleteTopic(string)

	// Subscribe 订阅
	Subscribe(topic, key string) bool

	// UnSubscribe 退订
	UnSubscribe(topic, key string) bool

	// GetTopicList 获取订阅topic的用户
	GetTopicList(topic string) []string
}
