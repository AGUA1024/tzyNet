/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"tzyNet/tCommon"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/net"
)

const (
	SevType_GateWay = iota
)

func main() {
	// 创建Room服务
	roomSev, err := tNet.NewService("0.0.0.0:8001", "Room")
	if err != nil {
		tCommon.Logger.SystemErrorLog(err)
	}

	// 设置路由解析器
	parser := tNet.NewPkgParser[*tNet.DefaultPbPkgParser]()
	roomSev.BindPkgParser(parser)

	// 路由注册
	net.RoomRouteRegister(roomSev)

	// 服务器启动
	roomSev.Start()
}
