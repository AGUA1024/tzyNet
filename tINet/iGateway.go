package tINet

import "net/http"

type IGateway interface{
	IService

	// 注册请求路径
	RoutePath(routePath string) IServerRouteGroup
	
	// 设置上线处理函数
	SetOnLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行上线处理函数
	GetOnLineHookFunc() func(ctx *tCommon.ConContext)

	// 设置断线处理函数
	SetOffLineHookFunc(fun func(ctx *tCommon.ConContext))
	// 执行断线处理函数
	GetOffLineHookFunc() func(ctx *tCommon.ConContext)

	// 服务器客户端连接初始化与注册
	ConRegister(respRw http.ResponseWriter, req *http.Request) (ICon, error)

	// 监听和处理客户端连接
	ListenAndHandle(con ICon)

	// 网络数据处理
	MsgHandle(req IRequest)
}