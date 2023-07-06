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
}
