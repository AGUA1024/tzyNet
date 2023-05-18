package common

import (
	"github.com/gorilla/websocket"
	"net/http"
	"path/filepath"
)

type gameMasterObj struct {
	conn *websocket.Conn
}

var ROOT_PATH, _ = filepath.Abs("./")
var CONFIG_PATH, _ = filepath.Abs("./config")
var PROTO_PATH, _ = filepath.Abs("./proto")
var LOG_PATH, _ = filepath.Abs("./log")
var API_PATH, _ = filepath.Abs("./api")

var GameMaster gameMasterObj

// 设置websocket, CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
