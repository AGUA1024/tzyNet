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
	conId      uint64
	conn       *websocket.Conn
	clientIp   string
	clientPort uint32
	property   map[string]any
}

func (this *WsCon) SetProperty(key string, value any) {

}

func (this *WsCon) GetProperty(key string) (any, error) {
	return nil, nil
}

func (this *WsCon) ReadMsg() (tINet.IMsg, error) {
	_, p, err := this.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	pkgParser := Service.GetPkgParser()
	return pkgParser.UnMarshal(p)
}

func (this *WsCon) GetRequest(msg tINet.IMsg) tINet.IRequest {
	return &wsRequest{
		msg: msg,
		con: this,
	}
}

// 添加连接
func (this WsConMaster) ConAdd(con tINet.ICon) {
	this.mpCon[con.GetConId()] = con
}

// 获取连接
func (this WsCon) GetCon() *websocket.Conn {
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
