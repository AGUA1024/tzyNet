package api

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/control"
	"hdyx/model"
)

func Act1GameInit(ctx *common.ConContext, params []byte) {
	outObj := control.Act1GameInit(ctx)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

func PlayCard(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.PlayCard_InObj](params, &api.PlayCard_InObj{})
	outObj := control.PlayCard(ctx, parmObj.CardIndex)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

func EventHandle(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.EventHandle_InObj](params, &api.EventHandle_InObj{})

	outObj := control.EventHandle(ctx, parmObj.ChooseIndex)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

// 摸牌
func GetCardFromPool(ctx *common.ConContext, params []byte) {
	outObj := control.GetCardFromPool(ctx)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

// 超时提醒
func TurnTimeOut(ctx *common.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.TimeOut_InObj](params, &api.TimeOut_InObj{})
	outObj := control.TurnTimeOut(ctx, parmObj.SeqId)

	if outObj != nil {
		model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
	}
}

func GetAct1Info(ctx *common.ConContext, params []byte) {
	outObj := control.GetAct1Info(ctx)

	model.MsgRoomBroadcast[*api.GetAct1Info_OutObj](ctx, outObj)
}
