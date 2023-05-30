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
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "创建失败，房间id已存在")
	} else if err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_QUERY, err)
	}

	// 创建房间
	if !model.CreateRoom(ctx, roomId, actId) {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_WRITE_ERROR, "创建房间时数据写入错误")
	}

	return &api.CreateRoom_OutObj{Ok: true}
}

// 销毁房间
func DestoryRoom(ctx *common.ConContext, roomId uint64) *api.DestroyRoom_OutObj {
	// 房间不存在
	roomModel, err := model.GetRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "要删除的房间不存在")
	}

	roomInfo := roomModel.UidToPlayerInfo
	// 是否在房间中
	if _, ok := roomInfo[ctx.GetConGlobalObj().Uid]; !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_IN_ROOM, "不在房间中，无法删除房间")
	}
	// 是否是房主
	if isMaster := roomInfo[ctx.GetConGlobalObj().Uid].IsMaster; !isMaster {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法删除房间")
	}

	// 销毁房间
	ok := model.DestroyRoom(ctx, roomId)
	if !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_WRITE_ERROR, "删除房间时，写入redis局部内存异常")
	}

	return &api.DestroyRoom_OutObj{Ok: true}
}

func JoinRoom(ctx *common.ConContext, roomId uint64) *api.JoinRoom_OutObj {
	// 房间不存在
	roomModel, err := model.GetRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "要加入的房间不存在")
	}

	roomInfo := roomModel.UidToPlayerInfo
	// 已经在房间中了
	if _, ok := roomInfo[ctx.GetConGlobalObj().Uid]; ok {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_ALREADY_IN_ROOM, "玩家已经在要加入的房间中了")
	}

	// 房间人已满

	// 加入房间
	roomInfo[ctx.GetConGlobalObj().Uid] = &model.PlayerModel{
		IsMaster: false,
		State:    false,
	}

	// 写入redis
	roomModel.UidToPlayerInfo = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = map[uint64]*api.PlayerInfo{}

	for uid, player := range roomInfo {
		mpRet[uid] = &api.PlayerInfo{
			State:    player.State,
			IsMaster: player.IsMaster,
		}
	}

	return &api.JoinRoom_OutObj{Players: mpRet}
}

func LeaveRoom(ctx *common.ConContext, roomId uint64) *api.LeaveRoom_OutObj {
	// 房间不存在
	roomModel, err := model.GetRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "游戏房间不存在,无法退出房间")
	}

	roomInfo := roomModel.UidToPlayerInfo

	// 是否在房间中
	if _, ok := roomInfo[ctx.GetConGlobalObj().Uid]; !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_IN_ROOM, "不在游戏房间中，无法退出房间")
	}

	// 是否是房主
	if isMaster := roomInfo[ctx.GetConGlobalObj().Uid].IsMaster; isMaster {
		// 房主转移
		for uid, player := range roomInfo {
			if uid == ctx.GetConGlobalObj().Uid {
				continue
			}

			// 设置新房主
			if player.IsMaster == false {
				roomInfo[uid].IsMaster = true
				break
			}
		}
	}

	// 离开房间
	delete(roomInfo, ctx.GetConGlobalObj().Uid)

	// 离开后如果没有玩家，则销毁房间
	if len(roomInfo) == 0 {
		ok := model.DestroyRoom(ctx, roomId)
		if !ok {
			common.Logger.GameErrorLog(ctx, common.ERR_REDIS_WRITE_ERROR, "删除房间时，写入redis局部内存异常")
		}
	}

	// 写入redis
	roomModel.UidToPlayerInfo = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = map[uint64]*api.PlayerInfo{}

	for uid, player := range roomInfo {
		mpRet[uid] = &api.PlayerInfo{
			State:    player.State,
			IsMaster: player.IsMaster,
		}
	}

	return &api.LeaveRoom_OutObj{Players: mpRet}
}

// 准备或取消准备
func SetOrCancelPrepare(ctx *common.ConContext, roomId uint64) *api.SetOrCancelPrepare_OutObj {
	// 房间不存在
	roomModel, err := model.GetRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "房间不存在，无法准备")
	}

	roomInfo := roomModel.UidToPlayerInfo

	// 是否在房间中
	if _, ok := roomInfo[ctx.GetConGlobalObj().Uid]; !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_IN_ROOM, "不在房间中，无法准备")
	}

	// 准备或取消准备
	roomInfo[ctx.GetConGlobalObj().Uid].State = !roomInfo[ctx.GetConGlobalObj().Uid].State

	// 写入redis
	roomModel.UidToPlayerInfo = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = map[uint64]*api.PlayerInfo{}

	for uid, player := range roomInfo {
		mpRet[uid] = &api.PlayerInfo{
			State:    player.State,
			IsMaster: player.IsMaster,
		}
	}

	return &api.SetOrCancelPrepare_OutObj{Players: mpRet}
}

func GameStart(ctx *common.ConContext, roomId uint64) *api.GameStart_OutObj {
	// 房间不存在
	roomModel, err := model.GetRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "游戏房间不存在，无法开启游戏")
	}

	roomInfo := roomModel.UidToPlayerInfo

	// 是否在房间中
	if _, ok := roomInfo[ctx.GetConGlobalObj().Uid]; !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_IN_ROOM, "不在房间中，无法开启游戏")
	}

	// 是否是房主
	if isMaster := roomInfo[ctx.GetConGlobalObj().Uid].IsMaster; !isMaster {
		common.Logger.GameErrorLog(ctx, common.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法开启游戏")
	}

	// 是否达到游戏开启人数
	//model := model.GetAct(ctx, roomModel.ActId)
	//model.(model.Act1Model)

	// 是否全部准备好开始游戏
	for _, player := range roomInfo {
		if player.State == false {
			common.Logger.GameErrorLog(ctx, common.ERR_IS_NO_PREPARE_EXIST, "玩家未全部准备好，无法开启游戏")
		}
	}

	return &api.GameStart_OutObj{Ok: true}
}
