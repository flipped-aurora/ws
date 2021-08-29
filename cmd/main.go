package main

import (
	"fmt"
	"os"

	"github.com/flipped-aurora/ws/core/biz"
	"github.com/flipped-aurora/ws/core/data"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"nhooyr.io/websocket"
)

// Init logger初始化
func Init() error {

	// 创建自定义logger
	var core zapcore.Core

	// 非发布版本
	devCore := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// 双输出
	core = zapcore.NewTee(
		// 添加到终端输出
		zapcore.NewCore(devCore, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
	)

	lg := zap.New(core, zap.AddCaller())
	// 替换zap库全局的logger对象
	// 使用zap.L() 调用全局对象
	zap.ReplaceGlobals(lg)

	// 启动哨兵
	return nil
}

func main() {
	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	m := data.NewManage(10)
	t := data.NewTopic()
	h := data.NewHandle()

	admin := data.NewAdmin(m, t, h, zap.L())
	admin.RegisteredMsgHandler(1, func(message biz.IMessage) bool {
		client, ok := admin.FindClient(message.GetTo())
		if !ok {
			admin.Logger.Info("没有找到该用户")
			return false
		}
		return client.SendMes(message)
	})
	adminCase := biz.NewAdmin(admin)
	r := gin.Default()
	r.GET("/ws", adminCase.HandlerWS(biz.CheckNull, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}))
	r.POST("/sendMeg", adminCase.SendMsg(biz.CheckNull))
	r.Run(":8090")
}
