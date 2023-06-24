/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"tzyNet/tCommon"
	"tzyNet/tNet"
)

func main() {
	flag.StringVar(&tNet.ENV_NODE_TAG, "tag", "", "Unique node label")
	flag.Parse()

	if len(tNet.ENV_NODE_TAG) == 0 {
		tCommon.Logger.SystemErrorLog("Invalid node tag")
	}

	tNet.ServerInit()
	wsServer := tNet.NewWsServer("0.0.0.0", 80, tNet.ENV_NODE_TAG)
	tNet.RouteRegister(wsServer)
	wsServer.Start()

	// 注册权限列表
	mpEtcdRegPerList, ok := tCommon.GetYamlMapCfg("serverCfg", "server", "release").(map[string]any)
	if !ok {
		tCommon.Logger.SystemErrorLog("ERR_ETCD_REGISTER_CFG: ")
	}

	// 判断是否有服务注册权限
	for _, reg := range mpEtcdRegPerList {
		strRegHost := reg.(string)

		// 有注册权限
		if tNet.ENV_NODE_HOST == strRegHost {
			// redis初始化
			rdb := redis.NewClient(&redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: "123456",
				DB:       0,
			})
			defer rdb.Close()
			// 使用rdb来执行Redis命令
			rdb.FlushAll(context.Background()).Result()
			break
		}
	}
}
