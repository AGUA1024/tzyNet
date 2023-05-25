package control

import (
	"context"
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
		common.Logger.GameErrorLog(ctx, common.ERR_NO_GATEWAY, "无可用网关")
	}

	// 网关选择
	nodeId := roomId % uint64(len(addrs))
	host := addrs[nodeId]

	// port
	resp, err := server.Etcd_Client.Get(context.Background(), server.Etcd_sevPort_Key)
	if err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_ETCD_GET_PORT, "从获取ETCD获取端口出错")
	}
	if len(resp.Kvs) == 0 {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_GATEWAY_PORT, "找不到网关端口")
	}

	port := string(resp.Kvs[0].Value)
	return &apiProto.GetGateWay_OutObj{
		Host: host,
		Port: port,
	}
}

// 注册用户信息
func ConGlobalObjInit(ctx *global.ConContext, cin *apiProto.ConGlobalObjInit_InObj) *apiProto.ConGlobalObjInit_OutObj {
	ok := ctx.SetConGlobalUid(cin.GetUid())
	if !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}

	if ok = ctx.SetConGlobalRoomId(cin.GetRoomId()); !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}
	return &apiProto.ConGlobalObjInit_OutObj{Ok: true}
}
