package data

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/flipped-aurora/ws/core/biz"
	"github.com/flipped-aurora/ws/core/utils"
)

// Manage 管理所有客户端
type Manage struct {
	// ctx: 上下文信息
	ctx context.Context

	// registry: 注册表
	registry [65536]slot

	// count
	count int64

	buf int64
}

type slot struct {
	// lock: 读写锁 避免并发
	sync.Mutex

	// map[string]biz.IClient
	// map[key]biz.IClient
	block map[string]biz.IClient
}

// Register 注册
func (m *Manage) Register(key string) biz.IClient {
	client := NewClient(m.ctx, m.buf)
	if v, ok := m.FindClient(key); ok {
		atomic.AddInt64(&m.count, -1)
		v.Shutdown()
	}
	shardingKey := utils.HashUint16(key)
	m.registry[shardingKey].Lock()
	defer m.registry[shardingKey].Unlock()
	if m.registry[shardingKey].block == nil {
		m.registry[shardingKey].block = make(map[string]biz.IClient)
	}
	m.registry[shardingKey].block[key] = client
	atomic.AddInt64(&m.count, 1)
	fmt.Println("注册:", key)
	return client
}

// UnRegister 注销
func (m *Manage) UnRegister(key string) {
	if _, ok := m.FindClient(key); !ok {
		return
	}
	shardingKey := utils.HashUint16(key)
	m.registry[shardingKey].Lock()

	if v, ok := m.registry[shardingKey].block[key]; ok {
		delete(m.registry[shardingKey].block, key)
		m.registry[shardingKey].Unlock()
		atomic.AddInt64(&m.count, -1)
		if v != nil && v != (biz.IClient)(nil) {
			v.Shutdown()
		}
	}
	fmt.Println("注销:", key)
}

// FindClient 查找客户端
func (m *Manage) FindClient(key string) (biz.IClient, bool) {
	shardingKey := utils.HashUint16(key)
	m.registry[shardingKey].Lock()
	defer m.registry[shardingKey].Unlock()
	v, ok := m.registry[shardingKey].block[key]
	return v, ok
}

// FindClients 批量查找客户端用户
func (m *Manage) FindClients(key ...string) []biz.IClient {
	res := make([]biz.IClient, 0)
	for _, s := range key {
		if v, ok := m.FindClient(s); ok {
			res = append(res, v)
		}
	}
	return res
}

// NewManage 新建 Manage对象
func NewManage(buf int64) *Manage {
	m := &Manage{
		buf: buf,
		ctx: context.Background(),
	}
	return m
}

// GetAll 查找所有客户端
func (m *Manage) GetAll() []biz.IClient {
	ret := make([]biz.IClient, 0, m.count)
	for index, _ := range m.registry {
		if m.registry[index].block == nil {
			continue
		}
		m.registry[index].Lock()
		for _, client := range m.registry[index].block {
			ret = append(ret, client)
		}
		m.registry[index].Unlock()
	}
	return ret
}
