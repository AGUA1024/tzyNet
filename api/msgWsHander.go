package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"hdyx/proto/act1"
	"net/http"
)

// 设置websocket
// CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func MsgWsHandler(c *gin.Context) {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close() //返回前关闭

	for {
		// 读取ws中的数据
		mt, _, err := ws.ReadMessage()
		if err != nil {
			break
		}

		remas := &act1.Student{
			Name:    "tzy",
			Age:     18,
			Address: "127.0.0.1",
			Cn:      1,
		}

		marsh, _ := proto.Marshal(remas)

		fmt.Println("protobuf:")
		fmt.Println(marsh)

		ms := &act1.Student{}
		proto.Unmarshal(marsh, ms)

		fmt.Println("name:", ms.Name)

		// 写入ws数据
		err = ws.WriteMessage(mt, marsh)
		if err != nil {
			break
		}
	}
}
