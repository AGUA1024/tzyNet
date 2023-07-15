package tINet

type IRequest interface {
	GetCon() ICon
	GetMsg() IMsg
	GetServName() string
}
