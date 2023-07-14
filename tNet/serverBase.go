package tNet

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"net"
	"strconv"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tServer"
)

const (
	Tcp uint16 = iota
	Http
	WebSocket
)

var Service tINet.IService = nil

type SevBase struct {
	host    string
	port    uint32
	sevName string
}

func (this *SevBase) GetHost() string {
	return this.host
}

func (this *SevBase) GetPort() uint32 {
	return this.port
}

func (this *SevBase) GetSevName() string {
	return this.sevName
}

func NewService(hostAddr string, protocolType uint16, sevName string) (tINet.IService, error) {
	switch protocolType {
	case Tcp:
	case Http:
		//return tServer.NewHttpServer(host, port, podName)
	case WebSocket:
		server, err := newWsService(hostAddr, sevName)
		if err != nil {
			return nil, err
		}
		return server, nil
	}

	return nil, errors.New("invalid_sev_type")
}

func ParseURL(urlStr string) (string, uint32, error) {
	host, port, err := net.SplitHostPort(urlStr)
	if err != nil {
		return "", 0, err
	}

	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "", 0, err
	}

	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return "", 0, err
	}

	return ipAddr.IP.String(), uint32(portNum), nil
}

// 服务器注册etcd
func etcdRegisterService(server tINet.IService) {
	// 创建租约
	cli := tServer.Etcd_Client
	sevName, ip, port := server.GetSevName(), server.GetHost(), server.GetPort()

	resp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", sevName, ip)
	value := fmt.Sprintf("%s:%d", ip, port)

	_, err = cli.Put(context.Background(), key, value, clientv3.WithLease(resp.ID))
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
