package tNet

import (
	"fmt"
	"reflect"
	"tzyNet/tCommon"
	"tzyNet/tNet/ioBuf"
	"tzyNet/tzyNet-demo/api"
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
func RouteRegister(server *WsServer) {
	routePath := server.RoutePath("/")
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

func (this *WsServer) getFuncByrouteCmd(cmd uint32) func(*tCommon.ConContext, []byte) {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff
	routeGroup, ok := this.routeGroup.mpRouteGroup[highCmd]
	if !ok {
		tCommon.Logger.SystemErrorLog("RouteGroup not found:", highCmd)
	}

	fun, ok := routeGroup.mpCmd[lowCmd]
	if !ok {
		fmt.Println(routeGroup.mpCmd)
		tCommon.Logger.SystemErrorLog("RouteCmd not found:", lowCmd)
	}
	return fun
}

// 路由处理
func wsRouteHandel(conCtx *tCommon.ConContext, cbuf *ioBuf.ClientBuf) {
	cmd := cbuf.CmdMerge
	byteApiBuf := cbuf.Data

	apiFunc := WsServerObj.getFuncByrouteCmd(cmd)

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {

		argValues := []reflect.Value{
			reflect.ValueOf(conCtx),
			reflect.ValueOf(byteApiBuf),
		}

		resultValues := fValue.Call(argValues)
		if len(resultValues) > 0 {
			//result := resultValues[0].Interface()
			//fmt.Println(result) // 输出：3
		} else {
			// 处理错误：函数没有返回值
		}
	} else {
		// 处理错误：f 不是一个函数类型
	}
}
