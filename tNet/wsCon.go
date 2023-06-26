package tNet

import (
	"github.com/gorilla/websocket"
)

type wsConnect struct {
	conn *websocket.Conn
}

func (master *wsConnect) GameDestroy() {
	// 关闭连接
	master.conn.Close()
}

func (master *wsConnect) OutMsg(messageType int, data []byte) error {
	return master.conn.WriteMessage(messageType, data)
}

func (master *wsConnect) ReadMsg() (messageType int, p []byte, err error) {
	return master.conn.ReadMessage()
}
