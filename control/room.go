package control

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/model"
)

// 创建房间
func CreateRoom(ctx *common.ConContext, roomId uint64, actId uint32) *api.CreateRoom_OutObj {
	// 房间已经存在
	ok, err := model.IsRoomExits(ctx, roomId)
	if ok {
		common.Logger.GameErrorLog(ctx, common.ERR_CREAT_ROOMID_EXIST, "创建失败，房间id已存在")
	} else if err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_QUERY, "Redis查询失败")
	}

	// 创建房间
	if model.CreateRoom(ctx, roomId, actId) {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_WRITE_ERROR, "创建房间时数据写入错误")
	}

	return &api.CreateRoom_OutObj{Ok: true}
}
