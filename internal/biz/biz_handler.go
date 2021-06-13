package biz

// TypeHandlerFunc 根据type处理handler
type TypeHandlerFunc func(IMessage) bool

// IHandle 注册对应type处理函数
type IHandle interface {
	Register(int32, TypeHandlerFunc) bool
	GetHandler(int32) (TypeHandlerFunc, bool)
}
