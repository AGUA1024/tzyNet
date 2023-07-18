package tINet

import "tzyNet/tIMiddleware"

type IService interface {
	// 启动服务
	Start()

	// 获取服务路由方法
	CallApiWithReq(request IRequest)

	// 绑定数据解析器
	BindPkgParser(IMsgParser)
	// 获取数据包解析器
	GetPkgParser() IMsgParser

	// 绑定Mq集群实例
	BindMq(mq tIMiddleware.IMq)
	// 获取Mq集群实例
	GetMq() tIMiddleware.IMq

	// 绑定Cache集群实例
	BindCache(cache tIMiddleware.ICache)
	// 获取Cache集群实例
	GetCache() tIMiddleware.ICache

	// 获取当前服务属性参数
	GetHost() string
	GetPort() uint32
	GetSevName() string
}

// 服务器路由组
type IServerRouteGroup interface {
	RouteGroup(routeGroup uint32) ISeverRoute
}

// 服务器路由
type ISeverRoute interface {
	Route(funcName uint32, request IRequest)
}
