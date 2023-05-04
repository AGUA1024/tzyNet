package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"hdyx/common"
	"hdyx/proto/ioBuf"
	"hdyx/server"
	"net/http"
	"regexp"
)

func anyRoute() {

}

var RouteRegister = map[int32]func(){
	1: anyRoute,
}

var R *gin.Engine = gin.Default()

// 设置websocket, CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	// 使用正则表达式匹配 URL，并提取请求路径
	re := regexp.MustCompile(`^ws?:\/\/[^\/]+([^?#]*)`)
	matches := re.FindStringSubmatch(server.ENV_GAME_URL)

	if len(matches) != 2 {
		common.Logger.ErrorLog("GAME_URL_ERROR")
	}

	reqPath := matches[1]

	fmt.Println(reqPath)
	R.GET(reqPath, ListenAndHandel)
}

func ListenAndHandel(c *gin.Context) {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.Logger.ErrorLog("WEBSOCKET_UPGRADE_ERROR")
	}

	defer ws.Close() //返回前关闭

	for {
		// 读取ws中的数据
		mt, msgBuf, err := ws.ReadMessage()
		if err != nil {
			break
		}

		clientBuf := &ioBuf.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil {
			common.Logger.ErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)

		outPut := &ioBuf.OutPutBuf{
			CmdCode:        12,
			ProtocolSwitch: 22,
			CmdMerge:       35,
			ResponseStatus: 67,
			ValidMsg:       "hello world",
			Data:           nil,
		}

		msg, _ := proto.Marshal(outPut)

		// 写入ws数据
		err = ws.WriteMessage(mt, msg)
		if err != nil {
			break
		}
	}
}
