package net

import (
	"tzyNet/tINet"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/api"
	"tzyNet/tzyNet-demo/model"
)

// Game服务路由组
const(
	route_Group_Act1 = iota
)

// 路由
const (
	route_Act1_Act1GameInit uint32 = iota + 1
	route_Act1_PlayCard
	route_Act1_EventHandle
	route_Act1_GetCardFromPool
	route_Act1_TurnTimeOut
	route_Act1_GetAct1Info
)

// 路由注册
func GameRouteRegister(routeMaster tINet.IService) {
	// 注册断线处理函数
	routeMaster.SetOffLineHookFunc(model.LoseConnectFunc)

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	routeMaster.BindPkgParser(parser)

	// 路由注册
	routePath := routeMaster.RoutePath("/")
	{
		// Act1
		routeGroup := routePath.RouteGroup(route_Group_Act1)
		{
			routeGroup.Route(route_Act1_Act1GameInit, api.Act1GameInit)
			routeGroup.Route(route_Act1_PlayCard, api.PlayCard)
			routeGroup.Route(route_Act1_EventHandle, api.EventHandle)
			routeGroup.Route(route_Act1_GetCardFromPool, api.GetCardFromPool)
			routeGroup.Route(route_Act1_TurnTimeOut, api.TurnTimeOut)
			routeGroup.Route(route_Act1_GetAct1Info, api.GetAct1Info)
		}
	}
}
