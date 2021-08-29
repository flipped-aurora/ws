package biz

// IManage 只提供用户管理存储
type IManage interface {
	// Register 注册接口
	Register(key string) IClient

	// UnRegister 注销接口
	UnRegister(key string)

	// FindClient 查找在线客户端
	FindClient(key string) (IClient, bool)

	// FindClients 批量查找在线客户端
	FindClients(key ...string) []IClient

	// GetAll 查找所有客户端
	GetAll() []IClient
}
