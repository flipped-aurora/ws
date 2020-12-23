package data

import (
	"Songzhibin/ws/internal/biz"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	registry atomic.Value

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
	oMap := m.registry.Load().(map[string]biz.IClient)
	nMap := make(map[string]biz.IClient, len(oMap)+1)
	for s, iClient := range oMap {
		nMap[s] = iClient
	}
	nMap[key] = client
	m.registry.Store(nMap)
	fmt.Println("注册:", key)
	return client
}

// UnRegister: 注销
func (m *Manage) UnRegister(key string) {
	if _, ok := m.FindClient(key); !ok {
		return
	}
	m.Lock()
	oldMap := m.registry.Load().(map[string]biz.IClient)
	nMap := make(map[string]biz.IClient, len(oldMap)-1)
	var i biz.IClient
	for s, iClient := range oldMap {
		if s == key {
			i = iClient
			continue
		}
		nMap[s] = iClient
	}
	m.Unlock()
	if i != nil && i != (biz.IClient)(nil) {
		i.Shutdown()
	}
	fmt.Println("注销:", key)
}

// FindClient: 查找客户端
func (m *Manage) FindClient(key string) (biz.IClient, bool) {
	vMap := m.registry.Load().(map[string]biz.IClient)
	v, ok := vMap[key]
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
	m.registry.Store(make(map[string]biz.IClient))
	return m
}
