package tNet

import "tzyNet/tINet"

type wsRequest struct {
	msg tINet.IMsg
	con tINet.ICon
}

func (this *wsRequest) GetMsg() tINet.IMsg {
	return this.msg
}

func (this *wsRequest) GetCon() tINet.ICon {
	return this.con
}
