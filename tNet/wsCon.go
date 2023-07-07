package tNet

import (
	"github.com/gorilla/websocket"
	"tzyNet/tINet"
)

type WsConMaster struct {
	wsUpgrader *websocket.Upgrader
	maxConId   uint64
	mpCon      map[uint64]tINet.ICon
}

type WsCon struct {
	conId   uint64
	conn    *websocket.Conn
	Storage any
}

// 添加连接
func (this WsConMaster) ConAdd(con tINet.ICon) {
	this.mpCon[con.GetConId()] = con
}

// 获取连接
func (this WsCon) GetCon() *websocket.Conn{
	return this.conn
}

// 获取连接id
func (this *WsCon) GetConId() uint64 {
	return this.conId
}

// 关闭连接
func (this *WsCon) Close() {
	// 关闭连接
	this.conn.Close()
}

// 数据输出
func (this *WsCon) OutMsg(messageType int, data []byte) error {
	return this.conn.WriteMessage(messageType, data)
}

// 数据读取
func (this *WsCon) ReadMsg() (messageType int, p []byte, err error) {
	return this.conn.ReadMessage()
}
