package control

import (
	"context"
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/server"
)

// 获取游戏服网关
func GetGateWay(ctx *common.ConContext, roomId uint64) *api.GetGateWay_OutObj {
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
	return &api.GetGateWay_OutObj{
		Host: host,
		Port: port,
	}
}

// 注册用户信息
func ConGlobalObjInit(ctx *common.ConContext, cin *api.ConGlobalObjInit_InObj) *api.ConGlobalObjInit_OutObj {
	ok := ctx.SetConGlobalUid(cin.GetUid())
	if !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}

	if ok = ctx.SetConGlobalRoomId(cin.GetRoomId()); !ok {
		common.Logger.GameErrorLog(ctx, common.ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}
	return &api.ConGlobalObjInit_OutObj{Ok: true}
}
