package biz

// Admin: 最终用户拿到的实体对象
type Admin struct {
	manage  IManage
	topic   ITopic
	handler IHandle
}

// NewAdmin: 实例化对象
func NewAdmin(manage IManage, topic ITopic, handler IHandle) *Admin {
	return &Admin{
		manage:  manage,
		topic:   topic,
		handler: handler,
	}
}
