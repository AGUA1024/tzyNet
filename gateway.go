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
	defer route.R.Run("localhost:8000")
}
