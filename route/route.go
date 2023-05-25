package route

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"hdyx/api"
	"hdyx/common"
	"hdyx/global"
	"hdyx/net/ioBuf"
	"reflect"
	"runtime"
	"strconv"
)

var GinEngine *gin.Engine = gin.Default()

func newConContext() *global.ConContext {
	return &global.ConContext{
		ConnectId: getGoroutineID(),
	}
}

// 注册connectGorutine全局空间
func registerConGlobal() *global.ConContext {
	ctx := newConContext()
	conId := ctx.GetConnectId()
	if _, ok := global.MpConRoutineStorage[conId]; ok {
		common.Logger.SystemErrorLog("RoutineStorage_MEM_OVERWRITE")
	}

	global.MpConRoutineStorage[conId] = global.RoutineStorage{
		Uid:    0,
		Cmd:    0,
		RoomId: 0,
	}

	return ctx
}

// 销毁connectGorutine全局空间
func destroyConGlobalObj(this *global.ConContext) {
	conId := this.GetConnectId()
	global.MpConRoutineStorage[conId].WsCon.Close()
	delete(global.MpConRoutineStorage, conId)
}

func getGoroutineID() uint64 {
	arrByte := make([]byte, 64)
	arrByte = arrByte[:runtime.Stack(arrByte, false)]
	arrByte = bytes.TrimPrefix(arrByte, []byte("goroutine "))
	arrByte = arrByte[:bytes.IndexByte(arrByte, ' ')]
	GoroutineID, _ := strconv.ParseUint(string(arrByte), 10, 64)
	return GoroutineID
}

func init() {
	GinEngine.GET("", ListenAndHandel)
}

func ListenAndHandel(ginCtx *gin.Context) {
	ws := common.WebSocketInit(ginCtx)

	// 注册connectGorutine全局空间
	conCtx := registerConGlobal()
	// 注册ws连接管道
	conCtx.SetConGlobalWsCon(ws)
	// 延迟注销connectGorutine全局空间,关闭ws连接
	defer destroyConGlobalObj(conCtx)

	for {
		// 读取ws中的数据
		_, msgBuf, err := ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(msgBuf)
		clientBuf := &ioBuf.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil {
			common.Logger.SystemErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					common.Logger.SystemErrorLog("PANIC_ERROR:", r)
				}
			}()

			routeHandel(conCtx, clientBuf)
		}()
	}
}

// 路由处理
func routeHandel(conCtx *global.ConContext, cbuf *ioBuf.ClientBuf) {
	cmd := cbuf.CmdMerge
	byteApiBuf := cbuf.Data
	apiFunc := api.GetApiByCmd(cmd)

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {

		argValues := []reflect.Value{
			reflect.ValueOf(conCtx),
			reflect.ValueOf(byteApiBuf),
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
