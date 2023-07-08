package tINet

type IPkgParser interface {
	NewParser() IPkgParser
	SetPkgObjBase()
	Marshal(obj any) ([]byte, error)
	UnMarshal([]byte) (IPkg, error)
}

type IPkg interface {
	GetRouteCmd() uint32
	GetData() []byte
	SetDataSrc(any)
}
