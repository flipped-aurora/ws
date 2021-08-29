package biz

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

const CheckNull = "CheckNull"
const RandomNum = 5

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type CheckFunc func(interface{}) (string, bool)

type AdminCase struct {
	IAdmin
}

// IAdmin 接口
type IAdmin interface {
	SendMsg(funcName string) func(c *gin.Context)
	HandlerWS(string, *websocket.AcceptOptions) func(c *gin.Context)
}

// NewAdmin 实例化对象
func NewAdmin(admin IAdmin) *AdminCase {
	return &AdminCase{
		IAdmin: admin,
	}
}

// RandStringBytesMaskImperSrc 生成随机字符串
func RandStringBytesMaskImperSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
