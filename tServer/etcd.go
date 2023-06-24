package tServer

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"tzyNet/tCommon"
)

var Etcd_Client *clientv3.Client
var Etcd_keepLive_ch <-chan *clientv3.LeaseKeepAliveResponse

// 服务器暴露接口
var Etcd_sevPort_Key = "sevPort"

func init() {
	endPoint := tCommon.GetYamlMapCfg("etcdCfg", "etcd", "host").(string)
	// 创建etcd客户端
	Etcd_Client, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{endPoint},
		DialTimeout: 5 * time.Second,
	})
	fmt.Println("--Etcd初始化完成")
}

// discoverService 从etcd中发现指定服务实例的地址
func DiscoverService(cli *clientv3.Client, serviceName string) ([]string, error) {
	// 查询服务实例
	resp, err := cli.Get(context.Background(), fmt.Sprintf("/services/%s", serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// 解析服务实例地址
	addrs := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}

	return addrs, nil
}
