/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"flag"
	"hdyx/common"
	"hdyx/net"
	"hdyx/server"
)

func main() {
	flag.StringVar(&server.ENV_NODE_TAG, "tag", "", "Unique node label")
	flag.Parse()

	if len(server.ENV_NODE_TAG) == 0 {
		common.Logger.SystemErrorLog("Invalid node tag")
	}

	server.ServerInit()
	defer net.GinEngine.Run("0.0.0.0:80")
}
