/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/net"
	"context"
)

func main() {
	// 创建webSokcet服务器对象
	wsServer := tNet.NewServer(tNet.SevType_WebSocketServer, "0.0.0.0", 8000, "node1")

	// 路由注册
	net.RouteRegister(wsServer)

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
