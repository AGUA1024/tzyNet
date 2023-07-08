package tNet

import (
	"fmt"
	"reflect"
	"tzyNet/tCommon"
	"tzyNet/tINet"
)

type RoutePathMaster struct {
	reqPath    string
	routeGroup *RouteGroupMaster // 路由组
	pkgParser  tINet.IPkgParser  // 路由解析器
}

type RouteGroupMaster struct {
	mpRouteGroup map[uint32]*RouteCmdMaster
}

type RouteCmdMaster struct {
	mpCmd map[uint32]func(ctx *tCommon.ConContext, params []byte)
}

// 绑定数据封包函数
func (this *RoutePathMaster) BindPkgParser(parser tINet.IPkgParser) {
	this.pkgParser = parser
}

// 将流数据转化为封包数据
func (this *RoutePathMaster) GetPkg(byteMsg []byte) (tINet.IPkg, error) {

	return this.pkgParser.UnMarshal(byteMsg)
}

// 请求路径
func (this *RoutePathMaster) RoutePath(routePath string) tINet.IServerRouteGroup {
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

func (this *WsServer) GetFuncByrouteCmd(cmd uint32) func(*tCommon.ConContext, []byte) {
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
func (this *WsServer) MsgHandle(conCtx *tCommon.ConContext, pkg tINet.IPkg) {
	cbuf := pkg.GetData()
	cmd := conCtx.GetConGlobalObj().Cmd

	apiFunc := Server.GetFuncByrouteCmd(cmd)

	fValue := reflect.ValueOf(apiFunc)
	if fValue.Kind() == reflect.Func {

		argValues := []reflect.Value{
			reflect.ValueOf(conCtx),
			reflect.ValueOf(cbuf),
		}

		resultValues := fValue.Call(argValues)
		if len(resultValues) > 0 {
			//result := resultValues[0].Interface()
			//fmt.Println(result) // 输出：3
		} else {
			// 处理错误：函数没有返回值
		}
	} else {
		// 处理错误：f 不是一个函数类型
	}
}
