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

	// 服务器客户端连接初始化
	ConRegister(respRw http.ResponseWriter, req *http.Request) ICon

	// 设置断线处理函数
	SetLoseConFunc(fun func(ctx *tCommon.ConContext))
	// 执行断线处理函数
	RunLoseConFunc(ctx *tCommon.ConContext)

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
