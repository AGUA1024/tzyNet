package tNet

import (
	"github.com/gorilla/websocket"
	"log"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tModel"
)

type WsConMaster struct {
	wsUpgrader *websocket.Upgrader
	maxConId   uint64
	mpCon      map[uint64]tINet.ICon
}

type WsCon struct {
	conId    uint64
	conn     *websocket.Conn
	property map[string]any
}

func (this *WsCon) SetProperty(key string, value any) {

}

func (this *WsCon) GetProperty(key string) (any, error) {
	return nil, nil
}

func (this WsCon) ListenAndHandle(conCtx *tCommon.ConContext) {
	con := this.conn
	for {
		// 读取ws中的数据
		msgType, msgBuf, err := con.ReadMessage()
		if err != nil {
			break
		}

		switch msgType {
		case websocket.TextMessage:
			// 处理文本消息
		case websocket.BinaryMessage:
			// 处理二进制消息
		case websocket.CloseMessage:
			log.Println("Received close message from client")
			// 处理关闭消息
			break
		}

		pkg, err := Server.GetPkg(msgBuf)
		if err != nil {
			tCommon.Logger.SystemErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		// 协程顺序执行
		done := make(chan bool, 1)

		go func() {
			defer func() {
				done <- true
				if r := recover(); r != nil {
					tCommon.Logger.SystemErrorLog("PANIC_ERROR:", r)
				}
			}()

			conCtx.EventStorageInit(pkg.GetRouteCmd())
			Server.MsgHandle(conCtx, pkg)

			tModel.AllRedisSave(conCtx)
		}()

		<-done
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

// 数据读取
func (this *WsCon) ReadMsg() (messageType int, p []byte, err error) {
	return this.conn.ReadMessage()
}
