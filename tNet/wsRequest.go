package tNet

import "tzyNet/tINet"

type wsRequest struct {
	sevName string
	msg tINet.IMsg
	con tINet.ICon
}

func (this *wsRequest) GetServName() string {
	return this.sevName
}

func (this *wsRequest) GetMsg() tINet.IMsg {
	return this.msg
}

func (this *wsRequest) GetCon() tINet.ICon {
	return this.con
}
