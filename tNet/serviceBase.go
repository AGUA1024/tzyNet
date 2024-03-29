package tNet

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.etcd.io/etcd/clientv3"
	"net"
	"net/http"
	"strconv"
	"tzyNet/tCommon"
	"tzyNet/tIMiddleware"
	"tzyNet/tINet"
	"tzyNet/tMiddleware/tServRegistry"
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
	mq      tIMiddleware.IMq
	cache   tIMiddleware.ICache
	db      tIMiddleware.IDb
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

func (this *SevBase) GetMq() tIMiddleware.IMq {
	return this.mq
}

func (this *SevBase) BindMq(mq tIMiddleware.IMq) {
	this.mq = mq
}

func NewGateWay(hostAddr string, protocolType uint16, sevName string) (tINet.IGateway, error) {
	switch protocolType {
	case Tcp:
	case Http:
		//return tServer.NewHttpGateway(host, port, podName)
	case WebSocket:
		server, err := newWsGateway(hostAddr, sevName)
		if err != nil {
			return nil, err
		}
		return server, nil
	}

	return nil, errors.New("invalid_sev_type")
}

func NewService(hostAddr string, sevName string) (tINet.IService, error) {
	ip, port, err := ParseURL(hostAddr)
	if err != nil {
		return nil, err
	}

	Service = &WsGateway{
		SevBase: SevBase{
			host:    ip,
			port:    port,
			sevName: sevName,
		},
		RouteMaster: RouteMaster{},
		ConMaster: &WsConMaster{
			wsUpgrader: &websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			mpCon:    make(map[uint64]tINet.ICon),
			maxConId: 0,
		},
	}
	return Service, nil
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
	cli := tServRegistry.Etcd_Client
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
	tServRegistry.Etcd_keepLive_ch, err = cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
	go func() {
		for range tServRegistry.Etcd_keepLive_ch {
		}
	}()

	fmt.Println("--服务注册完成")
}
