package tNet

import (
	"fmt"
	"tzyNet/tCommon"
	"tzyNet/tINet"
)

type RouteMaster struct {
	reqPath    string
	routeGroup *RouteGroupMaster // 路由组
	pkgParser  tINet.IMsgParser  // 路由解析器
}

type RouteGroupMaster struct {
	mpRouteGroup map[uint32]*RouteCmdMaster
}

type RouteCmdMaster struct {
	mpCmd map[uint32]func(ctx *tCommon.ConContext, params []byte)
}

// 绑定数据封包函数
func (this *RouteMaster) BindPkgParser(parser tINet.IMsgParser) {
	this.pkgParser = parser
}

// 将流数据转化为封包数据
func (this *RouteMaster) GetPkg(byteMsg []byte) (tINet.IMsg, error) {

	return this.pkgParser.UnMarshal(byteMsg)
}

// 请求路径
func (this *RouteMaster) RoutePath(routePath string) tINet.IServerRouteGroup {
	// 设置请求路径
	this.reqPath = routePath

	if this.routeGroup == nil {
		this.routeGroup = &RouteGroupMaster{
			mpRouteGroup: make(map[uint32]*RouteCmdMaster),
		}
	}

	return this.routeGroup
}

// 路由组
func (this *RouteGroupMaster) RouteGroup(routeGid uint32) tINet.ISeverRoute {
	if this.mpRouteGroup == nil {
		this.mpRouteGroup = make(map[uint32]*RouteCmdMaster)
	}
	_, ok := this.mpRouteGroup[routeGid]
	if !ok {
		this.mpRouteGroup[routeGid] = &RouteCmdMaster{
			mpCmd: make(map[uint32]func(ctx *tCommon.ConContext, params []byte)),
		}
	}
	return this.mpRouteGroup[routeGid]
}

// 路由注册
func (this *RouteCmdMaster) Route(cmd uint32, fun func(ctx *tCommon.ConContext, params []byte)) {
	_, ok := this.mpCmd[cmd]
	if ok {
		tCommon.Logger.SystemErrorLog("Duplicate routing configuration")
	}
	this.mpCmd[cmd] = fun
}

func (this *WsService) GetFuncByrouteCmd(cmd uint32) func(*tCommon.ConContext, []byte) {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff
	routeGroup, ok := this.routeGroup.mpRouteGroup[highCmd]
	if !ok {
		tCommon.Logger.SystemErrorLog("RouteGroup not found:", highCmd)
	}

	fun, ok := routeGroup.mpCmd[lowCmd]
	if !ok {
		fmt.Println(routeGroup.mpCmd)
		tCommon.Logger.SystemErrorLog("RouteCmd not found:", lowCmd)
	}
	return fun
}

// 路由处理
func (this *WsService) MsgHandle(msg tINet.IMsg) {
	this.mq.PushMsg(msg.GetServName(), msg)
}
