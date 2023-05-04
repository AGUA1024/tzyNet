package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"hdyx/proto/ioBuf"
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

	buf := "hello world"
	dataBuf := []byte(buf)
	// 发送消息
	message := &ioBuf.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       123,
		Data:           dataBuf,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}
	err = conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		panic(err)
	}

	// 接收消息
	_, resp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	receivedMessage := &ioBuf.OutPutBuf{}
	err = proto.Unmarshal(resp, receivedMessage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received message: %+v", receivedMessage)
}
