package data

import (
	"net/http"
	"time"

	"github.com/flipped-aurora/ws/core/biz"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

// Admin 最终用户拿到的实体对象
type Admin struct {
	biz.IManage
	biz.ITopic
	biz.IHandle
	Logger *zap.Logger
	// checkMap: 提供 k:v 插件方式注册校验函数等...
	// 使用者此处需要自行断言处理
	checkMap map[string]biz.CheckFunc
}

func NewAdmin(m biz.IManage, t biz.ITopic, h biz.IHandle, z *zap.Logger) *Admin {
	return &Admin{
		IManage:  m,
		ITopic:   t,
		IHandle:  h,
		Logger:   z,
		checkMap: make(map[string]biz.CheckFunc),
	}
}

// SendMsg 发送消息
func (a *Admin) SendMsg(funcName string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 校验
		key, ok := a.CheckWs(funcName, c)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": "身份验证失败",
			})
			return
		}
		// 先判断参数
		var msgBase Message
		if err := c.ShouldBind(&msgBase); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		msgBase.From = key
		msgBase.Time = time.Now().Unix()
		c.JSON(http.StatusOK, gin.H{
			"isOk": a.HandlerMes(&msgBase),
		})
		return
	}
}

// HandlerWS 注册web路由
func (a *Admin) HandlerWS(funcName string, options *websocket.AcceptOptions) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 升级连接
		conn, err := websocket.Accept(c.Writer, c.Request, options)
		if err != nil {
			a.Logger.Info("ws升级失败:", zap.Any("err", err))
			c.Status(101)
			return
		}
		// 校验
		key, ok := a.CheckWs(funcName, c)
		if !ok {
			a.Logger.Info("身份校验失败:")
			_ = conn.Close(websocket.StatusInvalidFramePayloadData, "身份校验失败")
			c.Status(401)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		client := a.Register(key)
		client.SetCtx(conn.CloseRead(client.GetCtx()))
		defer a.UnRegister(key)
		for {
			select {
			case <-client.GetCtx().Done():
				c.Status(499)
				return
			case msg, ok := <-client.MsgChan():
				if !ok {
					return
				}
				data, err := msg.Marshal()
				if err != nil {
					a.Logger.Error("Marshal error", zap.Any("err", err))
					return
				}
				_ = conn.Write(client.GetCtx(), websocket.MessageText, data)
			}
		}
	}
}

// ====================== 添加一些贫血方法 ========================

// Register 注册
func (a *Admin) Register(key string) biz.IClient {
	return a.IManage.Register(key)
}

// UnRegister 注销
func (a *Admin) UnRegister(key string) {
	a.IManage.UnRegister(key)
}

// AddCheckFunc 添加校验方法
func (a *Admin) AddCheckFunc(funcName string, f func(interface{}) (string, bool)) {
	if _, ok := a.checkMap[funcName]; ok {
		panic("重复添加function")
	}
	a.checkMap[funcName] = f
}

// CheckWs ws校验
func (a *Admin) CheckWs(funcName string, option interface{}) (string, bool) {
	// 如果和预留定义的不校验一致 直接返回true
	if funcName == biz.CheckNull {
		return biz.RandStringBytesMaskImperSrc(biz.RandomNum), true
	}
	// 判断注册的函数 如果没有注册 直接返回false
	f, ok := a.checkMap[funcName]
	if !ok {
		return "", false
	}
	return f(option)
}

// RegisteredMsgHandler 注册处理消息
func (a *Admin) RegisteredMsgHandler(t int32, handlerFunc biz.TypeHandlerFunc) bool {
	return a.IHandle.Register(t, handlerFunc)
}

// HandlerMes 处理message 分发客户端流程
func (a *Admin) HandlerMes(msg biz.IMessage) bool {
	f, ok := a.IHandle.GetHandler(msg.GetType())
	if !ok {
		return false
	}
	return f(msg)
}
