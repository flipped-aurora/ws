package data

import (
	"Songzhibin/ws/internal/biz"
	"sync"
	"sync/atomic"
)

type Topic struct {
	// lock: 读写锁 避免并发
	sync.Mutex

	// registry: 注册表
	// map[string]map[string]struct{}
	// map[topic]map[key]struct{}
	t atomic.Value
}

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

func (t *Topic) Publish(topic string, msg biz.IMessage) {
	oldTopic := t.t.Load().(map[string]map[string]struct{})
	if _, ok := oldTopic[topic]; !ok {
		return
	}
	// todo 未做嵌入manage
	//for key, _ := range oldTopic[topic] {
	//	client, ok := t.FindClient(key)
	//	if !ok {
	//		continue
	//	}
	//	client.SendMes(msg)
	//}
}
