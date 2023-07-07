package tINet

import (
	"net/http"
	"tzyNet/tCommon"
)

type IServer interface {
	// 启动服务器
	Start()

	// 注册服务器请求路径
	RoutePath(routePath string) IServerRouteGroup

	// 获取服务器路由方法
	GetFuncByrouteCmd(cmd uint32) func(*tCommon.ConContext, []byte)

	// 服务器客户端连接初始化与注册
	ConRegister(respRw http.ResponseWriter, req *http.Request) ICon

	// 设置上线处理函数
	SetOnLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行上线处理函数
	GetOnLineHookFunc() func(ctx *tCommon.ConContext)

	// 设置断线处理函数
	SetOffLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行断线处理函数
	GetOffLineHookFunc() func(ctx *tCommon.ConContext)

	// 绑定数据封包函数
	BindPkgParser(IPkgParser)
	// 将流数据转化为封包数据
	GetPkg(byteMsg []byte) (IPkg, error)

	// 获取服务器属性参数
	GetHost() string
	GetPort() uint32
	GetPodName() string
}

// 服务器路由组
type IServerRouteGroup interface {
	RouteGroup(routeGid uint32) ISeverRoute
}

// 服务器路由
type ISeverRoute interface {
	Route(cmd uint32, fun func(ctx *tCommon.ConContext, params []byte))
}
