package net

import (
	"tzyNet/tINet"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/api"
	"tzyNet/tzyNet-demo/model"
)

// Room路由组
const (
	route_Group_Room = iota
)

// 路由
const (
	route_Room_IsRoomExist uint32 = iota + 1
	route_Room_CreateRoom
	route_Room_JoinRoom
	route_Room_JoinGame
	route_Room_ChangePos
	route_Room_LeaveGame
	route_Room_SetOrCancelPrepare
	route_Room_GameStart
	route_Room_AddRobot
	route_Room_DelRobot
	route_Room_GetRoomInfo
)

// 路由注册
func RoomRouteRegister(routeMaster tINet.IService) {
	// 注册断线处理函数
	routeMaster.SetOffLineHookFunc(model.LoseConnectFunc)

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	routeMaster.BindPkgParser(parser)

	// 路由注册
	routePath := routeMaster.RoutePath("/")
	{
		// Room
		routeGroup := routePath.RouteGroup(route_Group_Room)
		{
			routeGroup.Route(route_Room_IsRoomExist, api.IsRoomExist)
			routeGroup.Route(route_Room_CreateRoom, api.CreateRoom)
			routeGroup.Route(route_Room_JoinRoom, api.JoinRoom)
			routeGroup.Route(route_Room_JoinGame, api.JoinGame)
			routeGroup.Route(route_Room_ChangePos, api.ChangePos)
			routeGroup.Route(route_Room_LeaveGame, api.LeaveGame)
			routeGroup.Route(route_Room_SetOrCancelPrepare, api.SetOrCancelPrepare)
			routeGroup.Route(route_Room_GameStart, api.GameStart)
			routeGroup.Route(route_Room_AddRobot, api.AddRobot)
			routeGroup.Route(route_Room_DelRobot, api.DelRobot)
			routeGroup.Route(route_Room_GetRoomInfo, api.GetRoomInfo)
		}
	}
}
