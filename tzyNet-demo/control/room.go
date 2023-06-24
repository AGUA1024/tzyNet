package control

import (
	"fmt"
	"strconv"
	api "tzyNet/tzyNet-demo/api/protobuf"
	"tzyNet/tzyNet-demo/model"
	"tzyNet/tzyNet-demo/sdk"
	"tzyNet/tCommon"
)

// 判断房间是否存在
func IsRoomExist(ctx *tCommon.ConContext, roomId uint64) *api.IsRoomExist_OutObj {
	// 房间已经存在
	ok, err := model.IsRoomExits(ctx, roomId)
	if err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_REDIS_QUERY, err)
	}

	return &api.IsRoomExist_OutObj{Ok: ok}
}

// 创建房间
func CreateRoom(ctx *tCommon.ConContext, roomId uint64, actId uint32, gameLv uint32) *api.CreateRoom_OutObj {
	uid := ctx.GetConGlobalObj().Uid
	// 房间已经存在
	ok, err := model.IsRoomExits(ctx, roomId)
	if ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "创建失败，房间id已存在")
	} else if err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_REDIS_QUERY, err)
	}

	// 获取用户信息
	arrUserInfo := sdk.GetUserInfoTestSdk.GetPlayerInfoBySdk(actId, uid)
	if arrUserInfo[uid] == nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_SDK_GET_USERINFO, "调取Sdk获取用户数据失败")
	}
	userInfo := arrUserInfo[uid]

	// 创建房间
	ret := model.CreateRoom(ctx, userInfo, roomId, actId, gameLv)
	if ret == nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_REDIS_WRITE_ERROR, "用户数据获取失败或创建房间时数据写入错误")
	}

	return ret
}

func JoinRoom(ctx *tCommon.ConContext, roomId uint64) *api.JoinRoom_OutObj {
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "要加入的房间不存在")
	}

	// 是否在房间中
	if ok, _ := model.GetRoomIndex(ctx, roomModel); ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_ALREADY_IN_ROOM, "已经在房间中，无法重复加入房间")
	}

	// 是否在游戏房中
	if ok, _ := model.GetGameIndex(ctx, roomModel); ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_ALREADY_IN_ROOM, "已经在房间中，无法重复加入房间")
	}

	// 当前房间信息
	roomInfo := roomModel.ArrUidAudience
	// 加入房间
	roomInfo = append(roomInfo, ctx.GetConGlobalObj().Uid)

	// 写入redis
	roomModel.ArrUidAudience = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(roomModel.PosIdToPlayer)
	return &api.JoinRoom_OutObj{
		Players: mpRet,
	}
}

func JoinGame(ctx *tCommon.ConContext, postionId uint32) *api.JoinGame_OutObj {
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "操作的房间不存在")
	}

	// 是否在观众席
	ok, roomIndex := model.GetRoomIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_ROOM, "不在房间中，无法加入游戏")
	}

	// 是否在游戏中
	ok, _ = model.GetGameIndex(ctx, roomModel)
	if ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_ALREADY_IN_GAME, "已经在游戏房间中了，无法重复加入")
	}

	// 游戏配置
	actCfg := model.GetActCfg(roomModel.ActId)
	// 游戏位置合法性
	if postionId > uint32(actCfg.MaxPlayerNum-1) {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_GAME_POSTION_PARAM, "错误的游戏位置参数:"+strconv.FormatUint(uint64(postionId), 10))
	}

	// 游戏房玩家数据
	gamePlayers := roomModel.PosIdToPlayer
	// 房间人已满
	if len(gamePlayers) == actCfg.MaxPlayerNum {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_ROOM_PLAYERNUM_FULL, "游戏房人已满，无法加入")
	}

	// 位次是否已经有人
	if _, ok = gamePlayers[postionId]; ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_ROOM_POSTION_USING, "新位次已经被玩家占用，无法调整游戏位次")
	}

	// 获取玩家信息
	arrUserInfo := sdk.GetUserInfoTestSdk.GetPlayerInfoBySdk(roomModel.ActId, ctx.GetConGlobalObj().Uid)
	userInfo, ok := arrUserInfo[ctx.GetConGlobalObj().Uid]
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_SDK_GET_USERINFO, "SDK获取玩家数据出错")
	}

	var isMaster = false
	// 游戏房间如果没人则成为房主
	playerNum := 0
	for _, player := range gamePlayers {
		if !player.IsRobot {
			playerNum++
		}
	}
	if playerNum == 0 {
		isMaster = true
	}

	player := &model.PlayerModel{
		Uid:      ctx.GetConGlobalObj().Uid,
		IsMaster: isMaster,
		State:    false,
		Head:     userInfo.Cover,
		Name:     userInfo.UserName,
	}

	// 加入房间
	gamePlayers[postionId] = player

	// 退出观众席
	roomModel.ArrUidAudience = append(roomModel.ArrUidAudience[:roomIndex], roomModel.ArrUidAudience[roomIndex+1:]...)

	// 写入redis
	roomModel.PosIdToPlayer = gamePlayers
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(gamePlayers)

	return &api.JoinGame_OutObj{Players: mpRet}
}

func ChangePos(ctx *tCommon.ConContext, newPosId uint32) *api.ChangePos_OutObj {
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "房间不存在，无法调整游戏位次")
	}

	// 是否在房间中
	ok, gameIndex := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_ROOM, "不在房间中，无法调整游戏位次")
	}

	// 位次入参合法性判断
	cfg := model.GetActCfg(roomModel.ActId)
	if newPosId >= uint32(cfg.MaxPlayerNum) {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_GAME_POSTION_PARAM, "游戏位次参数错误:"+strconv.FormatUint(uint64(gameIndex), 10))
	}

	roomInfo := roomModel.PosIdToPlayer
	// 新位次是否已经有人
	if _, ok = roomInfo[newPosId]; ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_ROOM_POSTION_USING, "新位次已经被玩家占用，无法调整游戏位次")
	}

	// 调整位次
	roomInfo[newPosId] = roomInfo[gameIndex]
	delete(roomInfo, gameIndex)

	// 写入redis
	roomModel.PosIdToPlayer = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(roomInfo)

	return &api.ChangePos_OutObj{Players: mpRet}
}

func LeaveGame(ctx *tCommon.ConContext) *api.LeaveGame_OutObj {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "游戏房间不存在,无法退出房间")
	}

	roomInfo := roomModel.PosIdToPlayer
	// 是否在房间
	ok, index := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_GAMEROOM, "不在游戏房间中，无法退出游戏")
	}

	if isMaster := roomInfo[index].IsMaster; isMaster {
		// 房主转移
		for i, player := range roomInfo {
			if i == index || player.IsRobot {
				continue
			}

			// 设置新房主
			if player.IsMaster == false {
				roomInfo[i].IsMaster = true
				break
			}
		}
	}

	// 离开房间
	delete(roomInfo, index)

	// 加入观众席
	roomModel.ArrUidAudience = append(roomModel.ArrUidAudience, ctx.GetConGlobalObj().Uid)

	// 写入数据
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(roomInfo)

	return &api.LeaveGame_OutObj{Players: mpRet}
}

// 准备或取消准备
func SetOrCancelPrepare(ctx *tCommon.ConContext) *api.SetOrCancelPrepare_OutObj {
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "房间不存在，无法准备")
	}

	// 是否在房间中
	ok, gameIndex := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_ROOM, "不在房间中，无法准备/取消准备")
	}

	roomInfo := roomModel.PosIdToPlayer
	// 准备或取消准备
	roomInfo[gameIndex].State = !roomInfo[gameIndex].State

	// 写入redis
	roomModel.PosIdToPlayer = roomInfo
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(roomInfo)

	return &api.SetOrCancelPrepare_OutObj{Players: mpRet}
}

func GameStart(ctx *tCommon.ConContext) *api.GameStart_OutObj {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "游戏房间不存在，无法开启游戏")
	}

	// 是否在游戏房间中
	ok, index := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_GAMEROOM, "不在游戏房间中，无法开始游戏")
	}

	// 是否是房主
	roomInfo := roomModel.PosIdToPlayer
	if isMaster := roomInfo[index].IsMaster; !isMaster {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法开启游戏")
	}

	//// 是否达到游戏开启人数
	//actCfg := model.GetActCfg(roomModel.ActId)
	//if len(roomInfo) < actCfg.MinPlayerNum {
	//	tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_PLAYERNUM_NOT_ENOUGH, "开启游戏人数不足")
	//}

	// 是否全部准备好开始游戏
	for _, player := range roomInfo {
		if player.State == false {
			tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NO_PREPARE_EXIST, "玩家未全部准备好，无法开启游戏")
		}
	}

	// 开始游戏
	roomModel.IsGame = true

	//写入数据
	model.RoomModelSave(ctx, roomModel)

	return &api.GameStart_OutObj{Ok: true}
}

func AddRobot(ctx *tCommon.ConContext, robotHead string, robotName string) *api.AddRobot_OutObj {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "游戏房间不存在，无法创建机器人")
	}

	// 是否在游戏房间中
	ok, index := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_GAMEROOM, "不在游戏房间中，无法创建机器人")
	}

	// 是否是房主
	roomInfo := roomModel.PosIdToPlayer
	if isMaster := roomInfo[index].IsMaster; !isMaster {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法创建机器人")
	}

	// 游戏配置
	actCfg := model.GetActCfg(roomModel.ActId)

	// 游戏房玩家数据
	gamePlayers := roomModel.PosIdToPlayer
	// 房间人已满
	if len(gamePlayers) == actCfg.MaxPlayerNum {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_ROOM_PLAYERNUM_FULL, "游戏房人已满，无法加入")
	}

	// 加入房间
	for pos := uint32(0); pos < uint32(actCfg.MaxPlayerNum); pos++ {
		if _, ok := gamePlayers[pos]; !ok {
			strRobotUid := fmt.Sprintf("%d%d", ctx.GetConGlobalObj().RoomId, pos)
			robotUid, _ := strconv.ParseUint(strRobotUid, 10, 64)

			player := &model.PlayerModel{
				Uid:      robotUid,
				IsMaster: false,
				State:    true,
				Head:     robotHead,
				Name:     robotName,
				IsRobot:  true,
			}

			gamePlayers[pos] = player
			break
		}
	}

	// 写入redis
	roomModel.PosIdToPlayer = gamePlayers
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	var mpRet = model.BackGameRoomInfo(gamePlayers)

	return &api.AddRobot_OutObj{Players: mpRet}
}

func DelRobot(ctx *tCommon.ConContext, robotPos uint32) *api.DelRobot_OutObj {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "游戏房间不存在，无法删除机器人")
	}

	// 是否在游戏房间中
	ok, index := model.GetGameIndex(ctx, roomModel)
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_GAMEROOM, "不在游戏房间中，无法删除机器人")
	}

	// 是否是房主
	roomInfo := roomModel.PosIdToPlayer
	if isMaster := roomInfo[index].IsMaster; !isMaster {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法删除机器人")
	}

	if _, ok = roomModel.PosIdToPlayer[robotPos]; !ok {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_PLAYER_ISNIT_EXIST, "该位置不是不存在玩家，无法执行删除机器人操作")
	}

	// 游戏房玩家数据
	if !roomModel.PosIdToPlayer[robotPos].IsRobot {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_PLAYER_ISNOT_ROBOT, "该位置不是机器人，无法删除")
	}

	// 踢出机器人
	delete(roomModel.PosIdToPlayer, robotPos)

	// 写入redis
	model.RoomModelSave(ctx, roomModel)

	// 返回数据
	gamePlayers := roomModel.PosIdToPlayer
	var mpRet = model.BackGameRoomInfo(gamePlayers)

	return &api.DelRobot_OutObj{Players: mpRet}
}

func GetRoomInfo(ctx *tCommon.ConContext) *api.GetRoomInfo_OutObj {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
	if roomModel == nil || err != nil {
		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "游戏房间不存在，无法创建机器人")
	}

	gamePlayers := roomModel.PosIdToPlayer
	var mpRet = model.BackGameRoomInfo(gamePlayers)

	return &api.GetRoomInfo_OutObj{Players: mpRet}
}

//// 销毁房间
//func DestoryRoom(ctx *tCommon.ConContext, roomId uint64) *api.DestroyRoom_OutObj {
//	// 房间不存在
//	roomModel, err := model.GetGameRoomInfo(ctx, roomId)
//	if roomModel == nil || err != nil {
//		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_NO_ROOMID_EXIST, "要删除的房间不存在")
//	}
//
//	// 是否在房间中
//	if ok := model.GetGameIndex(ctx, roomModel); !ok {
//		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_IN_ROOM, "不在房间中，无法删除房间")
//	}
//
//	// 是否是房主
//	if isMaster := roomInfo[ctx.GetConGlobalObj().Uid].IsMaster; !isMaster {
//		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_IS_NOT_ROOM_MASATER, "不是房主，无法删除房间")
//	}
//
//	// 销毁房间
//	ok := model.DestroyRoom(ctx, roomId)
//	if !ok {
//		tCommon.Logger.GameErrorLog(ctx, tCommon.ERR_REDIS_WRITE_ERROR, "删除房间时，写入redis局部内存异常")
//	}
//
//	return &api.DestroyRoom_OutObj{Ok: true}
//}
