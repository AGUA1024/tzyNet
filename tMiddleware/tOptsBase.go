package tMiddleware

type OptsBase interface{
	// 配置集群主机地址
	SetClusterHosts(hosts []string)

	// 获取集群主机地址
	GetClusterHosts() []string
}