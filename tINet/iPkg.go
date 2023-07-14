package tINet

type IMsgParser interface {
	NewParser() IMsgParser
	SetPkgObjBase()
	Marshal(obj any) ([]byte, error)
	UnMarshal([]byte) (IMsg, error)
}

type IMsg interface {
	GetRouteCmd() uint32
	GetData() []byte
	SetDataSrc(any)
}
