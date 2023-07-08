package tNet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tModel"
)

const (
	SevType_TcpServer uint16 = iota
	SevType_HttpServer
	SevType_WebSocketServer
)

type WsServer struct {
	OnLineFunc func(ctx *tCommon.ConContext)
	OffLineFun func(ctx *tCommon.ConContext)
	SeverBase
	RoutePathMaster
	ConMaster *WsConMaster
	pkgParser tINet.IPkgParser
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
			mpCon:    make(map[uint64]tINet.ICon),
			maxConId: 0,
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

		// 注册connectGorutine全局空间
		conCtx := tCommon.RegisterConGlobal(con.GetConId())

		// 延迟注销connectGorutine全局空间,关闭ws连接
		defer destroyConGlobalObj(conCtx)

		con.ListenAndHandle(conCtx)
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

// 绑定数据封包函数
func (this *WsServer) BindPkgParser(parser tINet.IPkgParser) {
	this.pkgParser = parser
}

// 将流数据转化为封包数据
func (this *WsServer) GetPkg(byteMsg []byte) (tINet.IPkg, error) {
	return this.pkgParser.UnMarshal(byteMsg)
}

// 设置上线处理函数
func (this *WsServer) SetOnLineHookFunc(fun func(ctx *tCommon.ConContext)) {
	this.OnLineFunc = fun
}

// 获取上线处理函数
func (this *WsServer) GetOnLineHookFunc() func(ctx *tCommon.ConContext) {
	return this.OnLineFunc
}

// 设置断线处理函数
func (this *WsServer) SetOffLineHookFunc(fun func(ctx *tCommon.ConContext)) {
	this.OffLineFun = fun
}

// 获取断线处理函数
func (this *WsServer) GetOffLineHookFunc() func(ctx *tCommon.ConContext) {
	return this.OffLineFun
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
	offLineFunc := Server.GetOffLineHookFunc()
	offLineFunc(conCtx)

	// cache数据落地
	tModel.AllRedisSave(conCtx)

	// 用户空间回收
	delete(tCommon.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(tCommon.MpConRoutineStorage, conCtx.ConnectId)
}
