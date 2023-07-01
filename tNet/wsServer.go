package tNet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync/atomic"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tModel"
	"tzyNet/tNet/ioBuf"
)

const (
	SevType_TcpServer uint16 = iota
	SevType_HttpServer
	SevType_WebSocketServer
)

type WsServer struct {
	Fun func(ctx *tCommon.ConContext)
	SeverBase
	RoutePathMaster
	ConMaster *WsConMaster
}

func newWsServer(host string, port uint32, podName string) tINet.IServer {
	Server = &WsServer{
		SeverBase: SeverBase{
			host:    host,
			port:    port,
			podName: podName,
		},
		RoutePathMaster: RoutePathMaster{},
		ConMaster: &WsConMaster{
			wsUpgrader: &websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			mpCon: make(map[uint64]tINet.ICon),
		},
	}
	return Server
}

// 启动服务器
func (this *WsServer) Start() {
	// 服务器服务注册
	etcdRegisterService(this)

	// 监听请求
	http.HandleFunc(this.reqPath, func(respRw http.ResponseWriter, req *http.Request) {
		con := this.ConRegister(respRw, req)
		defer con.Close()
		wsCon := con.GetCon()

		// 注册connectGorutine全局空间
		conCtx := tCommon.RegisterConGlobal(con.GetConId())
		// 注册ws连接管道
		conCtx.SetConGlobalWsCon(wsCon)
		// 延迟注销connectGorutine全局空间,关闭ws连接
		defer destroyConGlobalObj(conCtx)

		for {
			fmt.Println("reading-------------------")
			// 读取ws中的数据
			_, msgBuf, err := wsCon.ReadMessage()
			if err != nil {
				break
			}

			fmt.Println(msgBuf)
			clientBuf := &ioBuf.ClientBuf{}

			if err = proto.Unmarshal(msgBuf, clientBuf); err != nil || len(msgBuf) == 0 {
				tCommon.Logger.SystemErrorLog("PROTO_UNMARSHAL_ERROR")
			}

			fmt.Println("clientBuf：", clientBuf)
			fmt.Println("cmd:", clientBuf.CmdMerge)
			// 协程顺序执行
			done := make(chan bool, 1)

			go func() {
				defer func() {
					done <- true
					if r := recover(); r != nil {
						tCommon.Logger.SystemErrorLog("PANIC_ERROR:", r)
					}
				}()

				conCtx.EventStorageInit(clientBuf.CmdMerge)

				RouteHandel(conCtx, clientBuf)

				fmt.Println("uid:", conCtx.GetConGlobalObj().Uid)

				tModel.AllRedisSave(conCtx)
			}()

			<-done
		}
	})

	// 持续监听客户端连接
	serverAddr := fmt.Sprintf("%s:%d", this.GetHost(), this.GetPort())
	fmt.Println("[tzyNet] Server started successfully.")
	fmt.Printf("[tzyNet] Listen:%s:%d%s\n", this.GetHost(), this.GetPort(), this.reqPath)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
}

// 设置断线处理函数
func (this *WsServer) SetLoseConFunc(fun func(ctx *tCommon.ConContext))  {
	this.Fun = fun
}

// 执行断线处理函数
func (this *WsServer) RunLoseConFunc(ctx *tCommon.ConContext) {
	this.Fun(ctx)
}

// 连接注册
func (this *WsServer) ConRegister(respRw http.ResponseWriter, req *http.Request) tINet.ICon {
	// 升级WebSocket通信管道
	wsCon, err := this.ConMaster.wsUpgrader.Upgrade(respRw, req, nil)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 生成连接Id
	connectId := atomic.AddUint64(&this.ConMaster.maxConId, 1)
	con := &WsCon{
		conId: connectId,
		conn:  wsCon,
	}

	// 注册连接
	this.ConMaster.ConAdd(con)

	tCommon.MpConRoutineStorage[connectId] = &tCommon.ConGlobalStorage{
		WsCon:        nil,
		Uid:          0,
		Cmd:          0,
		RoomId:       0,
		EventStorage: nil,
	}
	return con
}

// 断开连接,退出游戏房间
func destroyConGlobalObj(conCtx *tCommon.ConContext) {
	// 如果在房间则离开房间，并广播
	Server.RunLoseConFunc(conCtx)

	// cache数据落地
	tModel.AllRedisSave(conCtx)

	// 用户空间回收
	delete(tCommon.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(tCommon.MpConRoutineStorage, conCtx.ConnectId)
}
