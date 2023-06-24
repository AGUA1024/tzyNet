package tNet

import (
	"tzyNet/tCommon"
)

type RouteMaster struct {
	reqPath    string
	routeGroup *RoutePathRegister // 路由组
}

type RoutePathRegister struct {
	mpRouteGroup map[uint32]*RouteCmdGroup
}

type RouteCmdGroup struct {
	mpCmd map[uint32]func(ctx *tCommon.ConContext, params []byte)
}

// 请求路径
func (this *RouteMaster) RoutePath(routePath string) *RoutePathRegister {
	// 设置请求路径
	this.reqPath = routePath

	if this.routeGroup == nil {
		this.routeGroup = &RoutePathRegister{
			mpRouteGroup: make(map[uint32]*RouteCmdGroup),
		}
	}

	return this.routeGroup
}

// 路由组
func (this *RoutePathRegister) RouteGroup(routeGid uint32) *RouteCmdGroup {
	if this.mpRouteGroup == nil {
		this.mpRouteGroup = make(map[uint32]*RouteCmdGroup)
	}
	_, ok := this.mpRouteGroup[routeGid]
	if !ok {
		this.mpRouteGroup[routeGid] = &RouteCmdGroup{
			mpCmd: make(map[uint32]func(ctx *tCommon.ConContext, params []byte)),
		}
	}
	return this.mpRouteGroup[routeGid]
}

// 路由注册
func (this *RouteCmdGroup) Route(cmd uint32, fun func(ctx *tCommon.ConContext, params []byte)) {
	_, ok := this.mpCmd[cmd]
	if ok {
		tCommon.Logger.SystemErrorLog("Duplicate routing configuration")
	}
	this.mpCmd[cmd] = fun
}
