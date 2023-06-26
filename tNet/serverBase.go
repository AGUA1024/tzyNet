package tNet

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tServer"
)

var Server tINet.IServer = nil

type SeverBase struct {
	host    string
	port    uint32
	podName string
}

func (this *SeverBase) GetHost() string {
	return this.host
}

func (this *SeverBase) GetPort() uint32 {
	return this.port
}

func (this *SeverBase) GetPodName() string {
	return this.podName
}

func NewServer(sevType uint16, host string, port uint32, podName string) tINet.IServer {
	switch sevType {
	case SevType_TcpServer:
	case SevType_HttpServer:
		//return tServer.NewHttpServer(host, port, podName)
	case SevType_WebSocketServer:
		return newWsServer(host, port, podName)
	}
	return nil
}

// 服务器注册etcd
func etcdRegisterService(server tINet.IServer) {
	// 创建租约
	cli := tServer.Etcd_Client
	podName,Host := server.GetPodName(), server.GetHost()

	resp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", podName, Host)
	_, err = cli.Put(context.Background(), key, Host, clientv3.WithLease(resp.ID))
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 设置心跳，避免租约过期
	tServer.Etcd_keepLive_ch, err = cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
	go func() {
		for range tServer.Etcd_keepLive_ch {
		}
	}()

	fmt.Println("--服务注册完成")
}
