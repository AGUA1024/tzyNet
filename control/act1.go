package control

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/model"
	"strconv"
)

func Act1GameInit(ctx *common.ConContext) *api.Act1Game_OutObj {
	// 获取房间信息
	roomInfo, err := model.GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if roomInfo == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "act1:操作的房间不存在")
	}

	// 房主是否开始游戏
	if roomInfo.IsGame == false {
		common.Logger.GameErrorLog(ctx, common.ERR_GAME_IS_NOT_STARTED, "act1:房主还未开始游戏")
	}

	// 游戏初始化
	act1model := model.NewActModel(ctx, 1)
	if act1model == nil {
		common.Logger.GameErrorLog(ctx, common.ERR_ACTID_IS_NOT_EXIST, "act1:获取model时参数错误，actId:"+strconv.Itoa(1))
	}

	outGameInfo := model.Act1InfoToOutObj(act1model.(*model.Act1Model))
	return &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: model.EVENT_TYPE_GAME_INIT,
		EventData: nil,
	}
}

func PlayCard(ctx *common.ConContext, cardIndex uint32) *api.Act1Game_OutObj {
	act1Model := getAct1Model(ctx)

	eventType, events := act1Model.PlayCard(ctx, int(cardIndex))
	if eventType == model.EVENT_TYPE_ERROR {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYCARD_PARAM, "act1:出牌时参数错误")
	}

	outGameInfo := model.Act1InfoToOutObj(act1Model)
	return &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: eventType,
		EventData: events,
	}
}

func EventHandle(ctx *common.ConContext, targetIndex uint32) *api.Act1Game_OutObj {
	act1Model := getAct1Model(ctx)

	// 权限判断
	evenPlayerIndex := act1Model.ActInfo.EventPlayer.PlayerIndex
	evenPlayer := act1Model.ActInfo.PlayerList[evenPlayerIndex]
	if evenPlayer.Uid != ctx.GetConGlobalObj().Uid {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_CANT_PLAY, "act1:没有出牌权限")
	}

	eventType, events := act1Model.EventHandler(ctx, targetIndex)
	if eventType == model.EVENT_TYPE_ERROR {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYCARD_PARAM, "act1:出牌时参数错误")
	}

	outGameInfo := model.Act1InfoToOutObj(act1Model)
	return &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: eventType,
		EventData: events,
	}
}

func GetCardFromPool(ctx *common.ConContext) *api.Act1Game_OutObj {
	act1Model := getAct1Model(ctx)
	act1Info := act1Model.ActInfo
	// 不是当前出牌人或者存在事件则报错
	if act1Info.BombPlayer != nil || ctx.GetConGlobalObj().Uid != act1Info.PlayerList[act1Info.CurPlayerIndex].Uid || act1Info.EventPlayer != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYCARD_PARAM, "act1:无法摸牌")
	}

	eventType, events := act1Model.GetCardFromPool(ctx)
	if eventType == model.EVENT_TYPE_ERROR {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYCARD_PARAM, "act1:无法摸牌")
	}

	outGameInfo := model.Act1InfoToOutObj(act1Model)
	return &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: eventType,
		EventData: events,
	}
}

// 获取对局信息,并判断请求合法性
func getAct1Model(ctx *common.ConContext) *model.Act1Model {
	// 获取房间信息
	roomInfo, err := model.GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if roomInfo == nil || err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_ROOMID_EXIST, "操作的房间不存在")
	}

	// 房主是否开始游戏
	if roomInfo.IsGame == false {
		common.Logger.GameErrorLog(ctx, common.ERR_GAME_IS_NOT_STARTED, "房主还未开始游戏")
	}

	// 是否可以出牌
	actModel := model.GetActModel(ctx, 1)
	if actModel == nil {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_GET_ACTINFO, "查询游戏信息失败")
	}

	// 获取对局信息
	act1Model := actModel.(*model.Act1Model)
	act1Info := act1Model.ActInfo

	// 是否已经出局
	if act1Info.PlayerList[act1Info.CurPlayerIndex].IsDie {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYER_IS_DIE, "已经出局，无法出牌")
	}

	return act1Model
}

// 超时提醒，报错不回包
func TurnTimeOut(ctx *common.ConContext, seqId uint32) *api.Act1Game_OutObj {
	act1Model := getAct1Model(ctx)
	act1Info := act1Model.ActInfo

	if act1Info.SeqId != seqId {
		return nil
	}

	eventType, events := act1Model.TurnTimeOut(ctx)

	if eventType == model.EVENT_TYPE_ERROR {
		common.Logger.GameErrorLog(ctx, common.ERR_ACT1_PLAYCARD_PARAM, "act1:出牌时参数错误")
	}

	outGameInfo := model.Act1InfoToOutObj(act1Model)
	return &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: eventType,
		EventData: events,
	}
}

func GetAct1Info(ctx *common.ConContext) *api.GetAct1Info_OutObj {
	actModel := model.GetActModel(ctx, 1)
	if actModel == nil {
		return &api.GetAct1Info_OutObj{
			GameInfo: nil,
		}
	}

	// 获取对局信息
	act1Model := actModel.(*model.Act1Model)

	// 是否需要断线重连
	if act1Model.IsPlayer(ctx) {
		act1Model.PlayerReConn(ctx)
	}

	outGameInfo := model.Act1InfoToOutObj(act1Model)
	return &api.GetAct1Info_OutObj{
		GameInfo: outGameInfo,
	}
}
