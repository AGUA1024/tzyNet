package control

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/global"
	"hdyx/model"
)

// 创建房间
func CreateRoom(ctx *global.ConContext, roomId uint64, actId uint32) *api.CreateRoom_OutObj {
	// 房间已经存在
	if model.RoomMaster.IsRoomExits(ctx, roomId) {
		common.Logger.GameErrorLog(ctx, common.ERR_CREAT_ROOMID_EXIST, "创建失败，房间id已存在")
	}

	// 创建房间
	model.RoomMaster.CreateRoom(ctx.GetConGlobalVal().Uid, roomId, actId)

	return &api.CreateRoom_OutObj{Ok: true}
}
