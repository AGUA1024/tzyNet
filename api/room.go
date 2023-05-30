package api

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/control"
	"hdyx/model"
)

// 创建房间
func CreateRoom(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.CreateRoom_InObj](params, &api.CreateRoom_InObj{})

	roomId := parmObj.GetRoomId()
	actId := parmObj.GetActId()

	outObj := control.CreateRoom(ctx, roomId, actId)

	common.OutPutStream[*api.CreateRoom_OutObj](ctx, outObj)
}

// 销毁房间
func DestroyRoom(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.DestroyRoom_InObj](params, &api.DestroyRoom_InObj{})

	outObj := control.DestoryRoom(ctx, parmObj.RoomId)

	model.MsgRoomBroadcast[*api.DestroyRoom_OutObj](ctx, outObj)
}

// 进入游戏房间
func JoinRoom(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.JoinRoom_InObj](params, &api.JoinRoom_InObj{})

	outObj := control.JoinRoom(ctx, parmObj.RoomId)

	model.MsgRoomBroadcast[*api.JoinRoom_OutObj](ctx, outObj)
}

// 离开游戏房间
func LeaveRoom(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.LeaveRoom_InObj](params, &api.LeaveRoom_InObj{})

	outObj := control.LeaveRoom(ctx, parmObj.RoomId)

	model.MsgRoomBroadcast[*api.LeaveRoom_OutObj](ctx, outObj)
}

// 准备
func SetOrCancelPrepare(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.SetOrCancelPrepare_InObj](params, &api.SetOrCancelPrepare_InObj{})

	outObj := control.SetOrCancelPrepare(ctx, parmObj.RoomId)

	model.MsgRoomBroadcast[*api.SetOrCancelPrepare_OutObj](ctx, outObj)
}

// 开始游戏
func GameStart(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.GameStart_InObj](params, &api.GameStart_InObj{})

	outObj := control.GameStart(ctx, parmObj.RoomId)

	model.MsgRoomBroadcast[*api.GameStart_OutObj](ctx, outObj)
}
