package common

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type wsConnect struct {
	conn *websocket.Conn
}

// 设置websocket, CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketInit(c *gin.Context) (*websocket.Conn, error) {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
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
