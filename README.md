# WebSocket

引用 nhooyr.io/websocket 进行二次封装 本工具包采用`CQRS`模式, 发送指令通过HTTP的接口,处理命令后通过WS通知WEB端

## 项目结构

``` shell
├── cmd
│ └── main.go  // 一个简单的Demo
└── internal 
    ├── biz // biz 定义了接口规范
    │ ├── biz_admin.go // 定义了IAdmin接口 IAdmin主要是定义了注册路由以及发送消息的接口
    │ ├── biz_client.go // 定义了IClient接口 IClient主要是定义了连接的ws的一些基础功能:发送消息、接收消息、获取ctx、设置ctx以及关闭
    │ ├── biz_handler.go // 定义了IHandle接口 IHandle主要是定义注册对应处理IMessage函数
    │ ├── biz_manage.go // 定义了IManage接口 IManage主要是定义ws得到管理存储提供了:注册、注销、查找在线客户端、批量查找在线、以及查找所有客户端
    │ ├── biz_message.go // 定义了IMessage接口 IMessage主要提供消息规范:编码、解码、获取消息类型、获取接收人类型
    │ └── biz_topic.go // 定义了ITopic接口 ITopic是扩展接口 可以指定一些topic去订阅对应的title例如:创建topic、删除topic、订阅、退订、获取topic的订阅用户
    ├── data // data 实现接口
    │ ├── data_admin.go // 实现了IAdmin接口
    │ ├── data_client.go // 实现了IClient接口
    │ ├── data_handler.go // 实现了IHandle接口
    │ ├── data_manage.go // 实现IManage接口
    │ ├── data_message.go // 实现了IMessage接口
    │ └── data_topic.go // 实现ITopic接口
    └── utils
        └── hash.go // 提供将string转化成一个uint16 用于一致性hash分slot
```

## 实例

这里用到 `cmd` 下的demo进行注解演示

```go
package main

import (
	"fmt"
	"os"
	"github.com/flipped-aurora/ws/internal/biz"
	"github.com/flipped-aurora/ws/internal/data"

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
	// 一些初始化...
	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建Manage对象指定下发消息的buffer容量
	m := data.NewManage(10)
	// 创建topic对象
	t := data.NewTopic()
	// 创建Handle路由组对象
	h := data.NewHandle()

	// admin 才是最终的管理对象, 将Manage、topic、Handle对象以及日志
	admin := data.NewAdmin(m, t, h, zap.L())
	// 注册处理消息对象 这里的意思是: 类型为1的就交到这个函数去处理
	admin.RegisteredMsgHandler(1, func(message biz.IMessage) bool {
		// 判断客户端是否存在
		client, ok := admin.FindClient(message.GetTo())
		if !ok {
			admin.Logger.Info("没有找到该用户")
			return false
		}
		// 存在转发消息
		return client.SendMes(message)
	})
	// 转化为接口对象 调用接口方法
	adminCase := biz.NewAdmin(admin)
	r := gin.Default()
	// biz.CheckNull 是跳过校验 否则的话是要 需要调用 admin.AddCheckFunc(name,func(interface)) 去进行校验
	// 比如jwt的话就可以嵌入进去做一个校验,校验成功后才可以连接,options则是ws的配置 自行选择
	// 多说一句, HandlerWS()第一个参数传定义好的funcname, CheckFunc func(interface{}) (string, bool) interface传入的其实是上下文,可以通过ctx来获取中间件保存的信息
	r.GET("/ws", adminCase.HandlerWS(biz.CheckNull, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}))
	// 发送消息 同样要走校验,校验通过发送消息 需要携带 Message 作为 post的参数
	r.POST("/sendMeg", adminCase.SendMsg(biz.CheckNull))
	r.Run(":8090")
}

```

## 其他

例子里没有用到的一些隐藏方法

- 如果使用了biz.CheckNull做funcName 他会调用 biz.RandStringBytesMaskImperSrc(biz.RandomNum)生成随机字符串
- data.Admin 除了实现 `IAdmin` 接口还有一些方法需要使用
  `AddCheckFunc(funcName string, f func(interface{})` 添加`SendMsg、HandlerWS`中传入funcName校验的方法
  `RegisteredMsgHandler(t int32, handlerFunc biz.TypeHandlerFunc`  注册处理消息的对应函数
- topic的注册之类的在demo中未体现,需要用户根据实际需求自行添加
    