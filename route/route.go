package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hdyx/common"
	ioBuf2 "hdyx/net/ioBuf"
	"hdyx/server"
	"net/http"
	"reflect"
	"regexp"
)

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

		clientBuf := &ioBuf2.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil {
			common.Logger.ErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)

		RouteHandel(clientBuf.CmdMerge, clientBuf.Data)

		outPut := &ioBuf2.OutPutBuf{
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

// 路由处理
func RouteHandel(cmd uint32, byteArg []byte) {
	api := GetApiByCmd(cmd)

	// 获取入参类型
	bufType := api.GetInType()

	param := reflect.New(bufType).Interface().(protoreflect.ProtoMessage)

	proto.Unmarshal(byteArg, param)

	apiFunc := api.GetFunc()

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {

		argValues := []reflect.Value{
			reflect.ValueOf(param),
		}

		resultValues := fValue.Call(argValues)
		if len(resultValues) > 0 {
			//result := resultValues[0].Interface()
			//fmt.Println(result) // 输出：3
		} else {
			// 处理错误：函数没有返回值
		}
	} else {
		// 处理错误：f 不是一个函数类型
	}
}

func GetInt32Cmd() {

}
