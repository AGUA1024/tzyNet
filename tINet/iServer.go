package tINet

import (
	"net/http"
	"tzyNet/tCommon"
)

type IService interface {
	// 启动服务
	Start()

	// 注册服务请求路径
	RoutePath(routePath string) IServerRouteGroup

	// 获取服务路由方法
	CallApiWithReq(request IRequest)

	GetFuncByrouteCmd(cmd uint32) func(*tCommon.ConContext, []byte)

	// 服务器客户端连接初始化与注册
	ConRegister(respRw http.ResponseWriter, req *http.Request) (ICon, error)

	// 设置上线处理函数
	SetOnLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行上线处理函数
	GetOnLineHookFunc() func(ctx *tCommon.ConContext)

	// 设置断线处理函数
	SetOffLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行断线处理函数
	GetOffLineHookFunc() func(ctx *tCommon.ConContext)

	// 绑定数据解析器
	BindPkgParser(IMsgParser)
	// 获取数据包解析器
	GetPkgParser() IMsgParser

	// 监听和处理客户端连接
	ListenAndHandle(con ICon)

	// 获取服务器属性参数
	GetHost() string
	GetPort() uint32
	GetSevName() string

	// 网络数据处理
	MsgHandle(req IRequest)
}

// 服务器路由组
type IServerRouteGroup interface {
	RouteGroup(routeGroup uint32) ISeverRoute
}

// 服务器路由
type ISeverRoute interface {
	Route(funcName uint32, request IRequest)
}
