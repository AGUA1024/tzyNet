package api

import (
	"fmt"
	api "tzyNet/tzyNet-demo/api/protobuf"
	"tzyNet/tzyNet-demo/control"
	"tzyNet/tzyNet-demo/model"
	"tzyNet/tCommon"
)

// 判断房间是否存在
func IsRoomExist(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.IsRoomExist_InObj](params, &api.IsRoomExist_InObj{})

	outObj := control.IsRoomExist(ctx, parmObj.RoomId)

	fmt.Println(outObj.Ok)
	tCommon.OutPutStream[*api.IsRoomExist_OutObj](ctx, outObj, tCommon.CONST_RESPONSE_STATUS_OK)
}

// 创建游戏房、观众房
func CreateRoom(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.CreateRoom_InObj](params, &api.CreateRoom_InObj{})

	roomId := parmObj.RoomId
	actId := parmObj.ActId
	gameLv := parmObj.GameLv
	outObj := control.CreateRoom(ctx, roomId, actId, gameLv)

	tCommon.OutPutStream[*api.CreateRoom_OutObj](ctx, outObj, tCommon.CONST_RESPONSE_STATUS_OK)
}

// 加入观众房
func JoinRoom(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.JoinRoom_InObj](params, &api.JoinRoom_InObj{})

	roomId := parmObj.RoomId
	outObj := control.JoinRoom(ctx, roomId)

	tCommon.OutPutStream[*api.JoinRoom_OutObj](ctx, outObj, tCommon.CONST_RESPONSE_STATUS_OK)
}

// 加入游戏位次
func JoinGame(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.JoinGame_InObj](params, &api.JoinGame_InObj{})

	posId := parmObj.PostionId
	outObj := control.JoinGame(ctx, posId)

	model.MsgRoomBroadcast[*api.JoinGame_OutObj](ctx, outObj)
}

// 更换到其他游戏位次
func ChangePos(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.ChangePos_InObj](params, &api.ChangePos_InObj{})

	newPosId := parmObj.NewPosId

	outObj := control.ChangePos(ctx, newPosId)

	model.MsgRoomBroadcast[*api.ChangePos_OutObj](ctx, outObj)
}

// 离开游戏房间
func LeaveGame(ctx *tCommon.ConContext, params []byte) {
	outObj := control.LeaveGame(ctx)

	model.MsgRoomBroadcast[*api.LeaveGame_OutObj](ctx, outObj)
}

// 准备
func SetOrCancelPrepare(ctx *tCommon.ConContext, params []byte) {
	outObj := control.SetOrCancelPrepare(ctx)

	model.MsgRoomBroadcast[*api.SetOrCancelPrepare_OutObj](ctx, outObj)
}

// 开始游戏
func GameStart(ctx *tCommon.ConContext, params []byte) {
	outObj := control.GameStart(ctx)

	model.MsgRoomBroadcast[*api.GameStart_OutObj](ctx, outObj)
}

// 添加机器人
func AddRobot(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.AddRobot_InObj](params, &api.AddRobot_InObj{})

	outObj := control.AddRobot(ctx, parmObj.GetRobotHead(), parmObj.GetRobotName())

	model.MsgRoomBroadcast[*api.AddRobot_OutObj](ctx, outObj)
}

// 删除机器人
func DelRobot(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.DelRobot_InObj](params, &api.DelRobot_InObj{})

	outObj := control.DelRobot(ctx, parmObj.PosId)

	model.MsgRoomBroadcast[*api.DelRobot_OutObj](ctx, outObj)
}

// 获取当前游戏房信息
func GetRoomInfo(ctx *tCommon.ConContext, params []byte) {
	outObj := control.GetRoomInfo(ctx)

	model.MsgRoomBroadcast[*api.GetRoomInfo_OutObj](ctx, outObj)
}

//// 销毁房间
//func DestroyRoom(ctx *tCommon.ConContext, params []byte) {
//	parmObj := tCommon.GetParamObj[*api.DestroyRoom_InObj](params, &api.DestroyRoom_InObj{})
//
//	outObj := control.DestoryRoom(ctx, parmObj.RoomId)
//
//	model.MsgRoomBroadcast[*api.DestroyRoom_OutObj](ctx, outObj)
//}
