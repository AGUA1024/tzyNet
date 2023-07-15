/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"tzyNet/tCommon"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/model"
	"tzyNet/tzyNet-demo/net"
)

const (
	SevType_GateWay = iota
)

func main() {
	// 创建webSokcet服务
	wsServer, err := tNet.NewGateWay("0.0.0.0:8000", tNet.WebSocket, "GateWay")
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 注册断线处理函数
	wsServer.SetOffLineHookFunc(model.LoseConnectFunc)

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	wsServer.BindPkgParser(parser)

	// 路由注册
	net.GameRouteRegister(wsServer)

	//redis数据初始化
	rdb := redis.NewClient(&redis.Options{
		Addr:     "123.207.49.151:6379",
		Password: "123456",
		DB:       0,
	})
	defer rdb.Close()
	// 使用rdb来执行Redis命令
	buf, err := rdb.FlushAll(context.Background()).Result()
	if len(buf) > 0 {
		fmt.Println("rdb.FlushAll:", buf)
	}
	if err != nil {
		fmt.Println("rdb.FlushAll_Err:", err)
	}

	// 服务器启动
	wsServer.Start()
}
