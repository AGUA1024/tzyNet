/**
 * @Author: tanzhenyu
 * @Date: 2023/04/24 14：40
 */

package main

import (
	_ "hdyx/common"
	"hdyx/server"
)

func main() {
	defer server.R.Run("localhost:8000")
}
