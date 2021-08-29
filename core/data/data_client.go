package data

import (
	"sync/atomic"

	"github.com/flipped-aurora/ws/core/biz"

	"context"
)

const closeFlag int32 = 1

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	msg    chan biz.IMessage

	// 标志位 判断该客户端是否关闭
	isClose int32
}

// SendMes 发送消息(非阻塞,如果已经满了则快速失败)
func (c *Client) SendMes(msg biz.IMessage) bool {
	if atomic.LoadInt32(&c.isClose) != 0 {
		return false
	}
	select {
	case c.msg <- msg:
		return true
	default:
		return false
	}
}

// MsgChan 返回channel
func (c *Client) MsgChan() <-chan biz.IMessage {
	return c.msg
}

// GetCtx 获取ctx
func (c *Client) GetCtx() context.Context {
	return c.ctx
}

// SetCtx 设置ctx
func (c *Client) SetCtx(ctx context.Context) {
	c.ctx = ctx
}

// Shutdown 关闭
func (c *Client) Shutdown() {
	if atomic.SwapInt32(&c.isClose, closeFlag) != 0 {
		return
	}
	close(c.msg)
	c.msg = nil
	c.cancel()
}

// NewClient 返回实例化
func NewClient(ctx context.Context, buf int64) *Client {
	client := &Client{msg: make(chan biz.IMessage, buf)}
	client.ctx, client.cancel = context.WithCancel(ctx)
	return client
}
