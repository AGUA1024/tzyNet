package net

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"hdyx/api"
	"hdyx/common"
	"hdyx/model"
	"hdyx/net/ioBuf"
	"log"
	"reflect"
)

var GinEngine *gin.Engine = gin.Default()

// 连接断开;空间回收
// 断开连接,退出游戏房间
func destroyConGlobalObj(conCtx *common.ConContext) {
	// 如果在房间则离开房间，并广播
	model.LeaveRoomAndBroadcast(conCtx)

	// cache数据落地
	model.AllRedisSave(conCtx)

	// 用户空间回收
	delete(common.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(common.MpConRoutineStorage, conCtx.ConnectId)
}

func init() {
	GinEngine.GET("", ListenAndHandel)
}

func ListenAndHandel(ginCtx *gin.Context) {
	ws, err := common.WebSocketInit(ginCtx)
	if err != nil {
		common.Logger.SystemErrorLog("WebSocketInit_ERROR")
	}

	// 注册connectGorutine全局空间
	conCtx := common.RegisterConGlobal()
	// 注册ws连接管道
	conCtx.SetConGlobalWsCon(ws)
	// 延迟注销connectGorutine全局空间,关闭ws连接
	defer destroyConGlobalObj(conCtx)

	for {
		fmt.Println("reading")
		// 读取ws中的数据
		_, msgBuf, err := ws.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(msgBuf)
		clientBuf := &ioBuf.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil || len(msgBuf) == 0 {
			common.Logger.SystemErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)

		// 协程顺序执行
		done := make(chan bool, 1)

		go func() {
			defer func() {
				log.Println("断开")
				done <- true
				if r := recover(); r != nil {
					common.Logger.SystemErrorLog("PANIC_ERROR:", r)
				}
			}()

			conCtx.EventStorageInit(clientBuf.CmdMerge)
			routeHandel(conCtx, clientBuf)

			fmt.Println("uid:")
			fmt.Println(conCtx.GetConGlobalObj().Uid)

			model.AllRedisSave(conCtx)
		}()

		<-done
	}
}

// 路由处理
func routeHandel(conCtx *common.ConContext, cbuf *ioBuf.ClientBuf) {
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
