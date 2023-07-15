package middleware

import (
	"github.com/nats-io/nats.go"
	"tzyNet/tCommon"
	"tzyNet/tINet"
	"tzyNet/tMiddleware/tMq"
)

func SevMqInit(sev tINet.IService) {
	// 获取nats配置
	opt := tMq.NatsOpts{Opts: nats.GetDefaultOptions()}
	// 设置nats集群主机地址
	opt.SetClusterHosts([]string{"123.207.49.151:8000"})
	// 创建Mq实例
	mq, err := tMq.NewMq[tMq.Nats](&opt)
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}
	// 绑定Mq实例
	sev.BindMq(mq)
}
