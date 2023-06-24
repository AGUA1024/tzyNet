package api

import (
	"tzyNet/tCommon"
	api "tzyNet/tzyNet-demo/api/protobuf"
	"tzyNet/tzyNet-demo/control"
	"tzyNet/tzyNet-demo/model"
)

func Act1GameInit(ctx *tCommon.ConContext, params []byte) {
	outObj := control.Act1GameInit(ctx)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

func PlayCard(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.PlayCard_InObj](params, &api.PlayCard_InObj{})
	outObj := control.PlayCard(ctx, parmObj.CardIndex)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

func EventHandle(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.EventHandle_InObj](params, &api.EventHandle_InObj{})

	outObj := control.EventHandle(ctx, parmObj.ChooseIndex)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

// 摸牌
func GetCardFromPool(ctx *tCommon.ConContext, params []byte) {
	outObj := control.GetCardFromPool(ctx)

	model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
}

// 超时提醒
func TurnTimeOut(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.TimeOut_InObj](params, &api.TimeOut_InObj{})
	outObj := control.TurnTimeOut(ctx, parmObj.SeqId)

	if outObj != nil {
		model.MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, outObj)
	}
}

func GetAct1Info(ctx *tCommon.ConContext, params []byte) {
	outObj := control.GetAct1Info(ctx)

	model.MsgRoomBroadcast[*api.GetAct1Info_OutObj](ctx, outObj)
}
