package tNet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"sync/atomic"
	"tzyNet/tCommon"
	"tzyNet/tIMiddleware"
	"tzyNet/tINet"
)

type WsGateway struct {
	OnLineFunc func(ctx *tCommon.ConContext)
	OffLineFun func(ctx *tCommon.ConContext)
	SevBase
	RouteMaster
	ConMaster *WsConMaster
	pkgParser tINet.IMsgParser
}

func (this *WsGateway) BindCache(cache tIMiddleware.ICache) {
	this.cache = cache
}

func (this *WsGateway) GetCache() tIMiddleware.ICache {
	return this.cache
}

func newWsGateway(hostAddr string, sevName string) (tINet.IGateway, error) {
	ip, port, err := ParseURL(hostAddr)
	if err != nil {
		return nil, err
	}

	Service = &WsGateway{
		SevBase: SevBase{
			host:    ip,
			port:    port,
			sevName: sevName,
		},
		RouteMaster: RouteMaster{},
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
	return Service, nil
}

// 启动服务
func (this *WsGateway) Start() {
	// 服务器初始化
	this.ServerInit()

	// 监听请求
	http.HandleFunc(this.reqPath, func(respRw http.ResponseWriter, req *http.Request) {
		con, err := this.ConRegister(respRw, req)
		if err != nil {

		}
		defer con.Close()

		this.ListenAndHandle(con)
	})

	// 持续监听客户端连接
	serverAddr := fmt.Sprintf("%s:%d", this.GetHost(), this.GetPort())
	fmt.Println("[tzyNet] Service started successfully.")
	fmt.Printf("[tzyNet] Listen:%s:%d%s\n", this.GetHost(), this.GetPort(), this.reqPath)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
}

// 服务器初始化
func (this *WsGateway) ServerInit() {
	// 服务器服务注册
	etcdRegisterService(this)
}

func (this *WsGateway) ListenAndHandle(con tINet.ICon) {
	for {
		msg, err := con.ReadMsg()
		if err != nil {

		}

		this.MsgHandle(msg)
	}
}

func (this *WsGateway) CallApiWithReq(req tINet.IRequest) {
	msg := req.GetMsg()
	cmd := msg.GetRouteCmd()

	apiFunc := Service.GetFuncByrouteCmd(cmd)

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {
		argValues := []reflect.Value{
			reflect.ValueOf(req),
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

// 绑定数据封包函数
func (this *WsGateway) BindPkgParser(parser tINet.IMsgParser) {
	this.pkgParser = parser
}

// 将流数据转化为封包数据
func (this *WsGateway) GetPkgParser() tINet.IMsgParser {
	return this.pkgParser
}

// 设置上线处理函数
func (this *WsGateway) SetOnLineHookFunc(fun func(ctx *tCommon.ConContext)) {
	this.OnLineFunc = fun
}

// 获取上线处理函数
func (this *WsGateway) GetOnLineHookFunc() func(ctx *tCommon.ConContext) {
	return this.OnLineFunc
}

// 设置断线处理函数
func (this *WsGateway) SetOffLineHookFunc(fun func(ctx *tCommon.ConContext)) {
	this.OffLineFun = fun
}

// 获取断线处理函数
func (this *WsGateway) GetOffLineHookFunc() func(ctx *tCommon.ConContext) {
	return this.OffLineFun
}

// 连接注册
func (this *WsGateway) ConRegister(respRw http.ResponseWriter, req *http.Request) (tINet.ICon, error) {
	// 升级WebSocket通信管道
	wsCon, err := this.ConMaster.wsUpgrader.Upgrade(respRw, req, nil)
	if err != nil {
		return nil, err
	}

	ip, port, err := ParseURL(req.RemoteAddr)
	if err != nil {
		return nil, err
	}

	// 生成连接Id
	connectId := atomic.AddUint64(&this.ConMaster.maxConId, 1)
	con := &WsCon{
		conId:      connectId,
		conn:       wsCon,
		clientIp:   ip,
		clientPort: port,
		property:   make(map[string]any),
	}

	// 注册连接
	this.ConMaster.ConAdd(con)

	MpConRoutineStorage[connectId] = make(map[string]any)

	return con, nil
}

// 断开连接,退出游戏房间
func destroyConGlobalObj(conCtx *tCommon.ConContext) {
	// 如果在房间则离开房间，并广播
	offLineFunc := Service.GetOffLineHookFunc()
	offLineFunc(conCtx)

	// cache数据落地
	tModel.AllRedisSave(conCtx)

	// 用户空间回收
	delete(tCommon.MpUserStorage, conCtx.ConnectId)

	// 关闭连接
	conCtx.GetConGlobalObj().WsCon.Close()

	// 销毁connectGorutine全局空间
	delete(MpConRoutineStorage, conCtx.ConnectId)
}
