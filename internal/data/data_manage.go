package data

import (
	"Songzhibin/ws/internal/biz"
	"context"
	"fmt"
	"sync"
)

// Manage 管理所有客户端
type Manage struct {
	// ctx: 上下文信息
	ctx context.Context
	// lock: 读写锁 避免并发
	sync.Mutex

	// registry: 注册表
	// map[string]biz.IClient
	// map[key]biz.IClient
	registry map[string]biz.IClient

	buf int64
}

// Register: 注册
func (m *Manage) Register(key string) biz.IClient {
	client := NewClient(m.ctx, m.buf)
	if v, ok := m.FindClient(key); ok {
		v.Shutdown()
	}
	m.Lock()
	defer m.Unlock()
	m.registry[key] = client
	fmt.Println("注册:", key)
	return client
}

// UnRegister: 注销
func (m *Manage) UnRegister(key string) {
	if _, ok := m.FindClient(key); !ok {
		return
	}
	m.Lock()
	if v, ok := m.registry[key]; ok {
		delete(m.registry, key)
		m.Unlock()
		if v != nil && v != (biz.IClient)(nil) {
			v.Shutdown()
		}
	}
	fmt.Println("注销:", key)
}

// FindClient: 查找客户端
func (m *Manage) FindClient(key string) (biz.IClient, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.registry[key]
	return v, ok
}

// FindClients: 批量查找客户端用户
func (m *Manage) FindClients(key ...string) []biz.IClient {
	res := make([]biz.IClient, 0)
	for _, s := range key {
		if v, ok := m.FindClient(s); ok {
			res = append(res, v)
		}
	}
	return res
}

func NewManage(buf int64) *Manage {
	m := &Manage{
		buf: buf,
		ctx: context.Background(),
	}
	m.registry = make(map[string]biz.IClient)
	return m
}

// GetAll: 查找所有客户端
func (m *Manage) GetAll() []biz.IClient {
	m.Lock()
	defer m.Unlock()
	res := make([]biz.IClient, 0, len(m.registry))
	for _, client := range m.registry {
		res = append(res, client)
	}
	return res
}
