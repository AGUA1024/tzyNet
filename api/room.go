package api

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/control"
	"hdyx/global"
)

// 创建房间
func CreateRoom(ctx *global.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.CreateRoom_InObj](params, &api.CreateRoom_InObj{})

	roomId := parmObj.GetRoomId()
	actId := parmObj.GetActId()

	outObj := control.CreateRoom(ctx, roomId, actId)

	common.OutPutStream[*api.CreateRoom_OutObj](ctx, outObj)
}

//// 销毁房间
//func DestroyRoom(uid uint64, params []byte) {
//
//}
//
//// 进入游戏房间
//func JoinRoom(uid uint64, params []byte) {
//	parmObj := GetParamObj[*api.JoinRoom_InObj](params, &api.JoinRoom_InObj{})
//	roomId := parmObj.GetRoomId()
//
//	outObj := &api.JoinRoom_OutObj{}
//	outStream := OutPutStream[*api.JoinRoom_OutObj](outObj)
//}
//
//// 离开游戏房间
//func LeaveRoom(uid uint64, params []byte) {
//	parmObj := GetParamObj[*api.LeaveRoom_InObj](params, &api.LeaveRoom_InObj{})
//
//	outObj := &api.LeaveRoom_OutObj{}
//
//	outStream := OutPutStream[*api.LeaveRoom_OutObj](outObj)
//}
//
//// 准备
//func RoomPrepare(uid uint64, params []byte) {
//	parmObj := GetParamObj[*api.LeaveRoom_InObj](params, &api.LeaveRoom_InObj{})
//
//	outObj := &api.LeaveRoom_OutObj{}
//
//	outStream := OutPutStream[*api.LeaveRoom_OutObj](outObj)
//}
//
//// 开始游戏
//func GameStart(uid uint64, params []byte) {
//	parmObj := GetParamObj[*api.LeaveRoom_InObj](params, &api.LeaveRoom_InObj{})
//
//	outObj := &api.LeaveRoom_OutObj{}
//
//	outStream := OutPutStream[*api.LeaveRoom_OutObj](outObj)
//}
