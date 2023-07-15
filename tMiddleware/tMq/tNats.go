package tMq

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"tzyNet/tIMiddleware"
	"tzyNet/tINet"
)

type Nats struct {
	con *nats.Conn
}

type NatsOpts struct {
	Opts nats.Options
}

func (this *NatsOpts) SetClusterHosts(hosts []string) {
	this.Opts.Servers = hosts
}

func (this *NatsOpts) GetClusterHosts() []string {
	return this.Opts.Servers
}

func (this *Nats) NewMq(opts tIMiddleware.IMqOpts) (tIMiddleware.IMq, error) {
	natsOpt, ok := opts.(*NatsOpts)
	if ok != true {
		return nil, errors.New("NewMq_Params_Must_Be_NatsOpts")
	}

	con, err := natsOpt.Opts.Connect()
	if err != nil {
		return nil, err
	}

	return &Nats{con: con}, nil
}

func (this *Nats) PushMsg(subject string, msg tINet.IMsg) error {
	return this.con.Publish(subject, msg.Serialize())
}

func (this *Nats) PopMsgWithCtx(subject string, ctx context.Context) ([]byte, error) {
	opts, err := this.con.SubscribeSync(subject)
	if err != nil {
		return nil, err
	}

	msg, err := opts.NextMsgWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
