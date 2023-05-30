package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	ioBuf2 "hdyx/net/ioBuf"
	"hdyx/script"
	"net/http"
	"os"
)

func main() {
	bufferShow()
	// 创建 HTTP 请求头
	header := http.Header{}
	//header.Add("Origin", "http://localhost")

	// 连接 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:80", header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chanel := make(chan string)

	go func() {
		for {
			select {
			case strBuf := <-chanel:
				data, err := hex.DecodeString(strBuf)
				err = conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					panic(err)
				}
				fmt.Println("write ok!")

				//
				//out := script.CreateRoom_OutObj{}
				//proto.Unmarshal(receivedMessage.Data, &out)
				//fmt.Println(out.Ok)
			}
		}
	}()

	go func() {
		for {
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
			fmt.Printf("Received message: %+v\n", receivedMessage)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		input := scanner.Text()

		if input == "quit" {
			break
		}

		chanel <- input
	}
}

func bufferShow() {
	roomId := uint64(6666)
	//_, str1 := script.GetGateWayBuf()
	//fmt.Println("GetGateWayBuf:", str1)

	_, str1 := script.ConGlobalObjInitBuf(1111, roomId)
	fmt.Println("ConGlobalObjInitBuf1:", str1)
	_, str11 := script.ConGlobalObjInitBuf(2222, roomId)
	fmt.Println("ConGlobalObjInitBuf11:", str11)
	_, str111 := script.ConGlobalObjInitBuf(333, roomId)
	fmt.Println("ConGlobalObjInitBuf111:", str111)

	_, str3 := script.CreateRoomBuf(roomId)
	fmt.Println("CreateRoomBuf：", str3)
	_, str4 := script.DestroyRoom(roomId)
	fmt.Println("DestroyRoom：", str4)
	_, str5 := script.JoinRoom(roomId)
	fmt.Println("JoinRoom:", str5)
	_, str6 := script.LeaveRoom(roomId)
	fmt.Println("LeaveRoom:", str6)
	_, str7 := script.SetOrCancelPrepare(roomId)
	fmt.Println("SetOrCancelPrepare:", str7)
	_, str8 := script.GameStart(roomId)
	fmt.Println("GameStart:", str8)
}
