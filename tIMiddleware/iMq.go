package tIMiddleware

import (
	"context"
	"tzyNet/tINet"
)

type IMq interface {
	NewMq(opts IMqOpts) (Mq IMq, err error)
	PushMsg(subject string, msg tINet.IMsg) error
	PopMsgWithCtx(subject string, ctx context.Context) ([]byte, error)
}

type IMqOpts interface {
	// 配置集群主机地址
	SetClusterHosts(hosts []string)
}
