package data

import (
	"sync"
	"sync/atomic"

	"github.com/flipped-aurora/ws/core/biz"
)

type Handle struct {
	// lock: 读写锁 避免并发
	sync.Mutex

	// f: 事件表
	// map[int32]biz.TypeHandlerFunc
	// map[type]biz.TypeHandlerFunc
	f atomic.Value
}

// Register 注册
func (h *Handle) Register(i int32, handlerFunc biz.TypeHandlerFunc) bool {
	h.Lock()
	defer h.Unlock()
	oMap := h.f.Load().(map[int32]biz.TypeHandlerFunc)
	if _, ok := oMap[i]; ok {
		panic("Repeat registration handlerFunc")
	}
	nMap := make(map[int32]biz.TypeHandlerFunc, len(oMap)+1)
	for k, v := range oMap {
		nMap[k] = v
	}
	nMap[i] = handlerFunc
	h.f.Store(nMap)
	return true
}

// GetHandler 获取注册函数
func (h *Handle) GetHandler(i int32) (biz.TypeHandlerFunc, bool) {
	oMap := h.f.Load().(map[int32]biz.TypeHandlerFunc)
	f, ok := oMap[i]
	return f, ok
}

func NewHandle() *Handle {
	h := &Handle{}
	h.f.Store(make(map[int32]biz.TypeHandlerFunc))
	return h
}
