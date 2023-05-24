/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	"flag"
	"hdyx/common"
	_ "hdyx/common"
	"hdyx/route"
	"hdyx/server"
)

func main() {
	flag.StringVar(&server.ENV_NODE_TAG, "tag", "", "Unique node label")
	flag.Parse()

	if len(server.ENV_NODE_TAG) == 0 {
		common.Logger.SystemErrorLog("Invalid node tag")
	}

	server.ServerInit()
	defer route.GinEngine.Run("0.0.0.0:80")
}
