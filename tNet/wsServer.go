package tNet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tNet/ioBuf"
	"tzyNet/tzyNet-demo/model"
)

const (
	SevType_TcpServer uint16 = iota
	SevType_HttpServer
	SevType_WebSocketServer
)

type WsServer struct {
	SeverBase
	RoutePathMaster
	Ws *websocket.Upgrader
}

var WsServerObj *WsServer = nil

func newWsServer(host string, port uint32, podName string) tINet.IServer {
	WsServerObj = &WsServer{
		SeverBase: SeverBase{
			host:    host,
			port:    port,
			podName: podName,
		},
		RoutePathMaster: RoutePathMaster{},
		Ws:              &websocket.Upgrader{},
	}

	return WsServerObj
}

func wsMsgHandler(respRw http.ResponseWriter, req *http.Request) {
	con, err := WsServerObj.Ws.Upgrade(respRw, req, nil)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
	defer con.Close()

	// 注册connectGorutine全局空间
	conCtx := tCommon.RegisterConGlobal()
	// 注册ws连接管道
	conCtx.SetConGlobalWsCon(con)
	// 延迟注销connectGorutine全局空间,关闭ws连接
	defer destroyConGlobalObj(conCtx)

	for {
		fmt.Println("reading")
		// 读取ws中的数据
		_, msgBuf, err := con.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(msgBuf)
		clientBuf := &ioBuf.ClientBuf{}

		if err = proto.Unmarshal(msgBuf, clientBuf); err != nil || len(msgBuf) == 0 {
			tCommon.Logger.SystemErrorLog("PROTO_UNMARSHAL_ERROR")
		}

		fmt.Println("clientBuf：")
		fmt.Println(clientBuf)
		fmt.Println("cmd:", clientBuf.CmdMerge)
		// 协程顺序执行
		done := make(chan bool, 1)

		go func() {
			defer func() {
				log.Println("断开")
				done <- true
				if r := recover(); r != nil {
					tCommon.Logger.SystemErrorLog("PANIC_ERROR:", r)
				}
			}()

			conCtx.EventStorageInit(clientBuf.CmdMerge)

			RouteHandel(conCtx, clientBuf)

			fmt.Println("uid:")
			fmt.Println(conCtx.GetConGlobalObj().Uid)

			model.AllRedisSave(conCtx)
		}()

		<-done
	}
}

func (this *WsServer) Start() {
	// 服务器服务注册
	etcdRegisterService(this)

	// 注册访问路径
	http.HandleFunc(this.reqPath, wsMsgHandler)

	// 持续监听客户端连接
	serverAddr := fmt.Sprintf("%s:%d", this.GetHost(), this.GetPort())
	fmt.Println("[tzyNet] Server started successfully.")
	fmt.Printf("[tzyNet] Listen:%s:%d%s\n", this.GetHost(), this.GetPort(), this.reqPath)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
}

// 连接断开;空间回收
// 断开连接,退出游戏房间
func destroyConGlobalObj(conCtx *tCommon.ConContext) {
	// 如果在房间则离开房间，并广播
	model.LeaveRoomAndBroadcast(conCtx)

	// cache数据落地
	model.AllRedisSave(conCtx)

	// 用户空间回收
	delete(tCommon.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(tCommon.MpConRoutineStorage, conCtx.ConnectId)
}
