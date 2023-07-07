package tINet

type IPkgParser interface {
	SetPkgObjBase(IPkg)
	Marshal([]byte, any) error
	UnMarshal([]byte) (IPkg, error)
}

type IPkg interface {
	GetRouteCmd() uint32
	GetData() []byte
}
