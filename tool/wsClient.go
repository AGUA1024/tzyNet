package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"hdyx/api/user"
	ioBuf2 "hdyx/net/ioBuf"
	"net/http"
)

func main() {
	// 创建 HTTP 请求头
	header := http.Header{}
	//header.Add("Origin", "http://localhost")

	// 连接 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8000/hdyx_game", header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := user.GetUserInfo_InObj{Uid: 1315}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x10001,
		Data:           cendBuf,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		panic(err)
	}

	// 接收消息
	_, resp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	receivedMessage := &ioBuf2.OutPutBuf{}
	err = proto.Unmarshal(resp, receivedMessage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received message: %+v", receivedMessage)
}
