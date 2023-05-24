package server

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hdyx/common"
	"io"
	"net/http"
	"runtime"
)

const (
	ENV_SERVER_PORT = 8000
	ENV_GAME_MARK   = "release" // 正式环境：release，开发环境:develop
	ENV_GAME_URL    = "ws://0.0.0.0:8000/hdyx_game"
)

var (
	ENV_NODE_TAG  = "" // 节点唯一标签
	ENV_NODE_HOST = ""
)

func ServerInit() {
	// 将 GOMAXPROCS 设置为 4
	runtime.GOMAXPROCS(4)

	// 获取节点公网ip
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		common.Logger.SystemErrorLog("ERROR_HOST_IP: ", err)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	ENV_NODE_HOST = string(ip)
	if err != nil || len(ENV_NODE_HOST) == 0 {
		common.Logger.SystemErrorLog("ERROR_HOST_IP: ", err)
		return
	}

	// 服务注册
	err = registerService(Etcd_Client, ENV_NODE_TAG, ENV_NODE_HOST, ENV_SERVER_PORT)
	if err != nil {
		common.Logger.SystemErrorLog("ETCD_REGISTER_ERROR: ", err)
	}
}

// registerService 向etcd注册服务信息
func registerService(cli *clientv3.Client, serviceName, host string, port int) error {
	// 创建租约
	resp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		return err
	}

	// 构建服务地址
	addr := fmt.Sprintf("%s:%d", host, port)

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", serviceName, addr)
	_, err = cli.Put(context.Background(), key, addr, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 设置心跳，避免租约过期
	etcd_keepLive_ch, err = cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	go func() {
		for range etcd_keepLive_ch {
		}
	}()

	return nil
}
