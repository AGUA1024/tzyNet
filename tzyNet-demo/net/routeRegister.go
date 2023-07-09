package net

import (
	"tzyNet/tINet"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/api"
	"tzyNet/tzyNet-demo/model"
)

// 路由组
const (
	routeGroup_GameNet uint32 = iota + 1
	routeGroup_Room
	routeGroup_Act1
)

// GameNet路由组
const (
	route_GameNet_GetGateWay uint32 = iota + 1
	route_GameNet_ConGlobalObjInit
)

// Room路由组
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

// Act1路由组
const (
	route_Act1_Act1GameInit uint32 = iota + 1
	route_Act1_PlayCard
	route_Act1_EventHandle
	route_Act1_GetCardFromPool
	route_Act1_TurnTimeOut
	route_Act1_GetAct1Info
)

// 路由注册
func RouteRegister(routeMaster tINet.IServer) {
	// 注册断线处理函数
	routeMaster.SetOffLineHookFunc(model.LoseConnectFunc)

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	routeMaster.BindPkgParser(parser)

	// 路由注册
	routePath := routeMaster.RoutePath("/")
	{
		// GameNet
		routeGroup := routePath.RouteGroup(routeGroup_GameNet)
		{
			routeGroup.Route(route_GameNet_GetGateWay, api.GetGateWay)
			routeGroup.Route(route_GameNet_ConGlobalObjInit, api.ConGlobalObjInit)
		}

		// Room
		routeGroup = routePath.RouteGroup(routeGroup_Room)
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

		// Act1
		routeGroup = routePath.RouteGroup(routeGroup_Act1)
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
