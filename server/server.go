package server

import "runtime"

const (
	ENV_SERVER_PORT = 8000
	ENV_GAME_MARK   = "release" // 正式环境：release，开发环境:develop
	ENV_GAME_URL    = "ws://0.0.0.0:8000/hdyx_game"

	VERSION = "1.0.0"
)

func init() {
	// 将 GOMAXPROCS 设置为 4
	runtime.GOMAXPROCS(4)
}
