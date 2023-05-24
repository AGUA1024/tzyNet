package control

import (
	apiProto "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/global"
	"hdyx/server"
)

// 获取游戏服网关
func GetGateWay(ctx *global.ConContext, roomId uint64) *apiProto.GetGateWay_OutObj {
	// 发现服务
	addrs, err := server.DiscoverService(server.Etcd_Client, "")
	if err != nil {
		common.Logger.GameErrorLog(ctx, "无可用网关")
	}

	// 网关选择
	nodeId := roomId % uint64(len(addrs))
	host := addrs[nodeId]

	return &apiProto.GetGateWay_OutObj{Host: host}
}

// 注册用户信息
func ConGlobalObjInit(ctx *global.ConContext, cin *apiProto.ConGlobalObjInit_InObj) *apiProto.ConGlobalObjInit_OutObj {
	ok := ctx.SetConGlobalUid(cin.GetUid())
	if !ok {
		common.Logger.GameErrorLog(ctx, "与服务器的连接不存在")
	}

	if ok = ctx.SetConGlobalRoomId(cin.GetRoomId()); !ok {
		common.Logger.GameErrorLog(ctx, "与服务器的连接不存在")
	}
	return &apiProto.ConGlobalObjInit_OutObj{Ok: true}
}
