package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"hdyx/api"
	"hdyx/common"
	ioBuf2 "hdyx/net/ioBuf"
	"hdyx/server"
	"reflect"
	"regexp"
)

var R *gin.Engine = gin.Default()

func init() {
	// 使用正则表达式匹配 URL，并提取请求路径
	re := regexp.MustCompile(`^ws?:\/\/[^\/]+([^?#]*)`)
	matches := re.FindStringSubmatch(server.ENV_GAME_URL)

	if len(matches) != 2 {
		common.Logger.ErrorLog("GAME_URL_ERROR")
	}

	//reqPath := matches[1]

	R.GET("", ListenAndHandel)
}

func ListenAndHandel(c *gin.Context) {
	common.GameMasterInit(c)
	defer common.GameMaster.GameDestroy()

	for {
		// 读取ws中的数据
		_, msgBuf, err := common.GameMaster.ReadMsg()
		if err != nil {
			break
		}

		clientBuf := &ioBuf2.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil {
			common.Logger.ErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					common.Logger.ErrorLog("PANIC_ERROR:", r)
				}
			}()
			routeHandel(clientBuf.CmdMerge, clientBuf.Data)
		}()
	}
}

// 路由处理
func routeHandel(cmd uint32, byteArg []byte) {
	apiFunc := api.GetApiByCmd(cmd)

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {

		argValues := []reflect.Value{
			reflect.ValueOf(byteArg),
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
