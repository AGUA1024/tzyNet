package control

import (
	"context"
	"tzyNet/tCommon"
	"tzyNet/tMiddleware/tServRegistry"
	"tzyNet/tzyNet-demo/app/gateWay/api/protobuf"
)

// 获取游戏服网关
func GetGateWay(ctx *tCommon.ConContext, roomId uint64) *api.GetGateWay_OutObj {
	// 发现服务
	addrs, err := tServRegistry.DiscoverService(tServRegistry.Etcd_Client, "")
	if err != nil {
		tCommon.Logger.GameErrorLog(ctx, ERR_NO_GATEWAY, "无可用网关")
	}

	// 网关选择
	nodeId := roomId % uint64(len(addrs))
	host := addrs[nodeId]

	// port
	resp, err := tServRegistry.Etcd_Client.Get(context.Background(), tServRegistry.Etcd_sevPort_Key)
	if err != nil {
		tCommon.Logger.GameErrorLog(ctx, ERR_ETCD_GET_PORT, "从获取ETCD获取端口出错")
	}
	if len(resp.Kvs) == 0 {
		tCommon.Logger.GameErrorLog(ctx, ERR_NO_GATEWAY_PORT, "找不到网关端口")
	}

	port := string(resp.Kvs[0].Value)
	return &api.GetGateWay_OutObj{
		Host: host,
		Port: port,
	}
}

// 注册用户信息
func ConGlobalObjInit(ctx *tCommon.ConContext, cin *api.ConGlobalObjInit_InObj) *api.ConGlobalObjInit_OutObj {
	ok := ctx.SetConGlobalUid(cin.GetUid())
	if !ok {
		tCommon.Logger.GameErrorLog(ctx, ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}

	if ok = ctx.SetConGlobalRoomId(cin.GetRoomId()); !ok {
		tCommon.Logger.GameErrorLog(ctx, ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}

	if !ctx.RegisterUserForStorage() {
		tCommon.Logger.GameErrorLog(ctx, ERR_NO_CONNECT_EXIST, "与服务器的连接不存在")
	}

	return &api.ConGlobalObjInit_OutObj{Ok: true}
}
