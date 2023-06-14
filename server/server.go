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

var (
	ENV_NODE_TAG  = "" // 节点唯一标签
	ENV_NODE_HOST = ""
)

func ServerInit() {
	// 将 GOMAXPROCS 设置为 4
	runtime.GOMAXPROCS(4)

	// 服务器注册etcd
	etcdRegisterService()
	fmt.Println("--服务注册完成")
}

// 服务器注册etcd
func etcdRegisterService() {
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

	// 注册权限列表
	mpEtcdRegPerList, ok := common.GetYamlMapCfg("serverCfg", "server", "release").(map[string]any)
	if !ok {
		common.Logger.SystemErrorLog("ERR_ETCD_REGISTER_CFG: ")
	}

	// 判断是否有服务注册权限
	for _, reg := range mpEtcdRegPerList {
		strRegHost := reg.(string)

		// 有注册权限
		if ENV_NODE_HOST == strRegHost {
			// 服务注册
			err = registerService(Etcd_Client, ENV_NODE_TAG, ENV_NODE_HOST)
			if err != nil {
				common.Logger.SystemErrorLog("ETCD_REGISTER_ERROR: ", err)
			}

			break
		}
	}
}

// registerService 向etcd注册服务信息
func registerService(cli *clientv3.Client, serviceName, host string) error {
	// 创建租约
	resp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		return err
	}

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", serviceName, host)
	_, err = cli.Put(context.Background(), key, host, clientv3.WithLease(resp.ID))
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
