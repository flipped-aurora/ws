/******
** @创建时间 : 2020/12/22 16:48
** @作者 : SongZhiBin
******/
package biz

// ITopic: 订阅表
type ITopic interface {
	// CreateTopic 创建topic
	CreateTopic(string)
	// DeleteTopic 删除topic
	DeleteTopic(string)

	// Subscribe 订阅
	Subscribe(topic, key string) bool

	// UnSubscribe 退订
	UnSubscribe(topic, key string) bool

	// Publish 发布
	Publish(topic string, msg IMessage)
}
