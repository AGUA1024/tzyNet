package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hdyx/api"
	"net/http"
)

// 设置websocket
// CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ping(c *gin.Context) {
	api.MsgWsHandler(c)
}
