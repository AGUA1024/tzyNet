package api

import (
	"tzyNet/tCommon"
	"tzyNet/tzyNet-demo/app/gateWay/api/protobuf"
	"tzyNet/tzyNet-demo/control"
)

// 获取游戏服网关
func GetGateWay(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.GetGateWay_InObj](params, &api.GetGateWay_InObj{})

	outObj := control.GetGateWay(ctx, parmObj.GetRoomId())

	tCommon.OutPutStream[*api.GetGateWay_OutObj](ctx, outObj, tCommon.CONST_RESPONSE_STATUS_OK)
}

// 连接全局变量初始化
func ConGlobalObjInit(ctx *tCommon.ConContext, params []byte) {
	parmObj := tCommon.GetParamObj[*api.ConGlobalObjInit_InObj](params, &api.ConGlobalObjInit_InObj{})

	outObj := control.ConGlobalObjInit(ctx, parmObj)

	tCommon.OutPutStream[*api.ConGlobalObjInit_OutObj](ctx, outObj, tCommon.CONST_RESPONSE_STATUS_OK)
}
