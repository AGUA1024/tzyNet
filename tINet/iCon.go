package tINet

import (
	"github.com/gorilla/websocket"
)

type ICon interface {
	// 获取连接Id
	GetConId() uint64

	// 获取与客户端建立的连接
	GetCon() *websocket.Conn

	// 关闭与客户端建立的连接
	Close()

	// 读取消息
	ReadMsg() (IMsg, error)

	// 获取完整的请求体
	GetRequest(msg IMsg) IRequest

	//SetProperty(key string, value any)
	//
	//GetProperty(key string) (any, error)
}
