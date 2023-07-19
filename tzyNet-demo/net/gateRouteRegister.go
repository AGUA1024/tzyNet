package net

import (
	"tzyNet/tINet"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/app/gateWay/api"
	"tzyNet/tzyNet-demo/model"
)

// GameNet路由组
const (
	route_Group_Gate = iota
)

// 路由
const (
	route_GameNet_GetGateWay uint32 = iota + 1
	route_GameNet_ConGlobalObjInit
)

// 路由注册
func GateRouteRegister(routeMaster tINet.IService) {
	// 注册断线处理函数
	routeMaster.SetOffLineHookFunc(model.LoseConnectFunc)

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	routeMaster.BindPkgParser(parser)

	// 路由注册
	routePath := routeMaster.RoutePath("/")
	{
		// GameNet
		routeGroup := routePath.RouteGroup(route_Group_Gate)
		{
			routeGroup.Route(route_GameNet_GetGateWay, api.GetGateWay)
			routeGroup.Route(route_GameNet_ConGlobalObjInit, api.ConGlobalObjInit)
		}
	}
}
