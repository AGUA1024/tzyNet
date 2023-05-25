package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	ioBuf2 "hdyx/net/ioBuf"
	"hdyx/script"
	"net/http"
)

func main() {
	// 创建 HTTP 请求头
	header := http.Header{}
	//header.Add("Origin", "http://localhost")

	// 连接 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8000", header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	data, str := script.ConGlobalObjInitBuf()
	fmt.Println(str)

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
	fmt.Printf("Received message: %+v\n", receivedMessage.CmdMerge)

	out := script.CreateRoom_OutObj{}
	proto.Unmarshal(receivedMessage.Data, &out)
	fmt.Println(out.Ok)
}
