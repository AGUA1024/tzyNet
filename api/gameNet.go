package api

import (
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/control"
	"hdyx/global"
)

// 获取游戏服网关
func GetGateWay(ctx *global.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.GetGateWay_InObj](params, &api.GetGateWay_InObj{})

	outObj := control.GetGateWay(ctx, parmObj.GetRoomId())

	common.OutPutStream[*api.GetGateWay_OutObj](ctx, outObj)
}

// 连接全局变量初始化
func ConGlobalObjInit(ctx *global.ConContext, params []byte) {
	parmObj := common.GetParamObj[*api.ConGlobalObjInit_InObj](params, &api.ConGlobalObjInit_InObj{})

	outObj := control.ConGlobalObjInit(ctx, parmObj)

	common.OutPutStream[*api.ConGlobalObjInit_OutObj](ctx, outObj)
}
