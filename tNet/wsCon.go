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

func (this WsConMaster) ConAdd(con tINet.ICon) {
	this.mpCon[con.GetConId()] = con
}

func (this WsCon) GetCon() *websocket.Conn{
	return this.conn
}

func (this *WsCon) GetConId() uint64 {
	return this.conId
}

func (this *WsCon) Close() {
	// 关闭连接
	this.conn.Close()
}

func (this *WsCon) OutMsg(messageType int, data []byte) error {
	return this.conn.WriteMessage(messageType, data)
}

func (this *WsCon) ReadMsg() (messageType int, p []byte, err error) {
	return this.conn.ReadMessage()
}
