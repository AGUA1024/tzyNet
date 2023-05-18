/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	_ "hdyx/common"
	"hdyx/route"
)

func main() {
	//uid := int32(1)
	//db := server.GetDb(uid, "hdyx_game")
	//db.InsertData(context.Background(), "act", info)
	defer route.R.Run("0.0.0.0:80")
}
