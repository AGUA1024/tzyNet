package tINet

import "tzyNet/tCommon"

type IServer interface {
	// 启动服务器
	Start()

	// 注册服务器请求路径
	RoutePath(routePath string) IServerRouteGroup

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
