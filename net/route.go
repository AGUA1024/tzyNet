package net

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"hdyx/api"
	"hdyx/common"
	"hdyx/control"
	"hdyx/model"
	"hdyx/net/ioBuf"
	"reflect"
	"runtime"
	"strconv"
)

var GinEngine *gin.Engine = gin.Default()

func newConContext() *common.ConContext {
	return &common.ConContext{
		ConnectId: getGoroutineID(),
	}
}

// 注册connectGorutine全局空间
func registerConGlobal() *common.ConContext {
	ctx := newConContext()
	conId := ctx.GetConnectId()
	if _, ok := common.MpConRoutineStorage[conId]; ok {
		common.Logger.SystemErrorLog("RoutineStorage_MEM_OVERWRITE")
	}

	common.MpConRoutineStorage[conId] = &common.ConGlobalStorage{
		WsCon:        nil,
		Uid:          0,
		Cmd:          0,
		RoomId:       0,
		EventStorage: nil,
	}

	return ctx
}

// 连接断开;空间回收
// 断开连接,退出游戏房间
func destroyConGlobalObj(conCtx *common.ConContext) {
	// 房间不存在
	roomModel, _ := model.GetRoomInfo(conCtx, conCtx.GetConGlobalObj().RoomId)
	if roomModel != nil {
		roomInfo := roomModel.UidToPlayerInfo
		// 是否在房间中
		if _, ok := roomInfo[conCtx.GetConGlobalObj().Uid]; ok {
			// 退出房间
			obj := control.LeaveRoom(conCtx, conCtx.GetConGlobalObj().RoomId)
			model.MsgRoomBroadcast(conCtx, obj)
		}
	}

	// cache数据落地
	model.AllRedisSave(conCtx)

	// 用户空间回收
	delete(common.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(common.MpConRoutineStorage, conCtx.ConnectId)
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
	ws, err := common.WebSocketInit(ginCtx)
	if err != nil {
		common.Logger.SystemErrorLog("WebSocketInit_ERROR")
	}

	// 注册connectGorutine全局空间
	conCtx := registerConGlobal()
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

		go func() {
			defer func() {
				if r := recover(); r != nil {
					common.Logger.SystemErrorLog("PANIC_ERROR:", r)
				}
			}()

			conCtx.EventStorageInit(clientBuf.CmdMerge)
			routeHandel(conCtx, clientBuf)

			fmt.Println("conGlobalObj:")
			fmt.Println(common.MpConRoutineStorage)
			fmt.Println(conCtx.GetConGlobalObj())
			//fmt.Println(conCtx.GetConGlobalObj().EventStorage)
			model.AllRedisSave(conCtx)
		}()
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
