package tMq

import "tzyNet/tIMiddleware"

func NewMq[mqType tIMiddleware.IMq](MqOpts tIMiddleware.IMqOpts) (tIMiddleware.IMq, error) {
	var mqObj mqType
	return mqObj.NewMq(MqOpts)
}
