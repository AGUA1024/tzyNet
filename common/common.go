package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"hdyx/global"
	"hdyx/net/ioBuf"
	"net/http"
	"os"
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

func WebSocketInit(c *gin.Context) *websocket.Conn {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		Logger.SystemErrorLog("WEBSOCKET_UPGRADE_ERROR")
	}

	return ws
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

func GetYamlMapCfg(fileName string, index ...string) any {
	//	返回值
	var ret any

	// 读取文件所有内容装到 []byte 中
	cfgFile := CONFIG_PATH + `/` + fileName + ".yaml"
	bytes, err := os.ReadFile(cfgFile)
	if err != nil {
		Logger.SystemErrorLog("GET_YAML_CFG_ERROR:" + cfgFile)
	}
	// 创建配置文件的结构体
	var conf map[string]any
	// 调用 Unmarshall 去解码文件内容
	// 注意要穿配置结构体的指针进去
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		Logger.SystemErrorLog("YAML_UNMARSHAL_ERROR:" + fmt.Sprintln(err))
	}

	for _, i := range index {
		var value any
		var ok bool

		if value, ok = conf[i]; !ok {
			Logger.SystemErrorLog("GET_YAML_CFG_INDEX_ERROR:" + fmt.Sprintln(err, "INDEX:", i))
		}

		// 是否读完
		if conf, ok = value.(map[string]any); ok {
			ret = conf
		} else {
			ret = value
		}
	}

	return ret
}

func GetParamObj[T proto.Message](params []byte, obj T) T {
	proto.Unmarshal(params, obj)
	return obj
}

func OutPutStream[T proto.Message](ctx *global.ConContext, obj T) {
	data, _ := proto.Marshal(obj)

	out := ioBuf.OutPutBuf{
		CmdCode:        0,
		ProtocolSwitch: 0,
		CmdMerge:       ctx.GetConGlobalVal().Cmd,
		ResponseStatus: 0,
		ValidMsg:       "",
		Data:           data,
	}

	outStream, err := proto.Marshal(&out)
	if err != nil {
		Logger.SystemErrorLog("GET_OUT_STREAM_Marshal_ERROR", err)
	}

	err = ctx.GetConGlobalVal().WsCon.WriteMessage(TextMessage, outStream)
	if err != nil {
		Logger.SystemErrorLog("OUT_STREAM_ERROR", err)
	}
}
