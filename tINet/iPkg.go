package tINet

type IMsgParser interface {
	NewParser() IMsgParser
	SetPkgObjBase()
	Marshal(obj any) ([]byte, error)
	UnMarshal([]byte) (IMsg, error)
}

type IMsg interface {
	// 获取biz服务名
	GetServName() string
	// 获取路由编码
	GetRouteCmd() uint32
	// 将结构体解析为二进制流数据
	Serialize() []byte
	// 二进制流数据反序列化为结构体数据
	Deserialize() IMsg

	SetDataSrc(any)
}
