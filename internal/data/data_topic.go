package data

import (
	"sync"
	"sync/atomic"
)

type Topic struct {
	// lock: 读写锁 避免并发
	sync.Mutex

	// t: 订阅表
	// map[string]map[string]struct{}
	// map[topic]map[key]struct{}
	t atomic.Value
}

func NewTopic() *Topic {
	t := &Topic{}
	t.t.Store(make(map[string]map[string]struct{}))
	return t
}

// CreateTopic 创建topic
func (t *Topic) CreateTopic(topic string) {
	t.Lock()
	defer t.Unlock()
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; ok {
		// 已经存在topic
		return
	}
	nTopic := make(map[string]map[string]struct{}, len(oldTopic)+1)
	for k, v := range oldTopic {
		nTopic[k] = v
	}
	nTopic[topic] = make(map[string]struct{})
	t.t.Store(nTopic)
}

// DeleteTopic 删除topic
func (t *Topic) DeleteTopic(topic string) {
	t.Lock()
	defer t.Unlock()
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; !ok {
		// 没有此topic 快速返回
		return
	}
	nTopic := make(map[string]map[string]struct{}, len(oldTopic)-1)
	for k, v := range oldTopic {
		if k == topic {
			continue
		}
		nTopic[k] = v
	}
	t.t.Store(nTopic)
}

// Subscribe 订阅
func (t *Topic) Subscribe(topic, key string) bool {
	t.Lock()
	defer t.Unlock()
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; !ok {
		return false
	}
	// 判断是否已经订阅
	if _, ok := oldTopic[topic][key]; ok {
		return false
	}
	oldTopic[topic][key] = struct{}{}
	return true
}

// UnSubscribe 退订
func (t *Topic) UnSubscribe(topic, key string) bool {
	t.Lock()
	defer t.Unlock()
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; !ok {
		return false
	}
	// 判断是否已经订阅
	if _, ok := oldTopic[topic][key]; !ok {
		return false
	}
	delete(oldTopic[topic], key)
	return true
}

// GetTopicList 获取订阅topic的用户
func (t *Topic) GetTopicList(topic string) []string {
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; !ok {
		// 表示没有此topic
		return nil
	}
	res := make([]string, 0)
	for key, _ := range oldTopic[topic] {
		res = append(res, key)
	}
	return res
}
