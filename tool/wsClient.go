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
<<<<<<< HEAD
=======

			fmt.Println("ResponseStatus:", receivedMessage.ResponseStatus)
			data := &script.Act1Game_OutObj{}
			proto.Unmarshal(receivedMessage.Data, data)

			fmt.Println("is:", data)
>>>>>>> origin/main
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
<<<<<<< HEAD
	roomId := uint64(6666)
	//_, str1 := script.GetGateWayBuf()
	//fmt.Println("GetGateWayBuf:", str1)
=======
	roomId := uint64(777)
	_, str0 := script.IsRoomExist(roomId)
	fmt.Println("IsRoomExist:", str0)
>>>>>>> origin/main

	_, str1 := script.ConGlobalObjInitBuf(13324782, roomId)
	fmt.Println("ConGlobalObjInitBuf1:", str1)
	_, str11 := script.ConGlobalObjInitBuf(11834295, roomId)
	fmt.Println("ConGlobalObjInitBuf11:", str11)
	_, str111 := script.ConGlobalObjInitBuf(11677176, roomId)
	fmt.Println("ConGlobalObjInitBuf111:", str111)

	_, str3 := script.CreateRoomBuf(roomId)
	fmt.Println("CreateRoomBuf：", str3)
	_, str4 := script.JoinRoom(roomId)
	fmt.Println("JoinRoom:", str4)

<<<<<<< HEAD
=======
	_, str24 := script.AddRobot("head1", "robot1")
	fmt.Println("AddRobot1:", str24)
	_, str25 := script.AddRobot("head2", "robot2")
	fmt.Println("AddRobot2:", str25)

	_, str27 := script.PlayerReCon()
	fmt.Println("PlayerReCon:", str27)

	_, str19 := script.DelRobot(1)
	fmt.Println("DelRobot1:", str19)
	_, str20 := script.DelRobot(2)
	fmt.Println("DelRobot2:", str20)
	_, str21 := script.DelRobot(3)
	fmt.Println("DelRobot3:", str21)
	_, str22 := script.DelRobot(4)
	fmt.Println("DelRobot4:", str22)
	_, str23 := script.DelRobot(5)
	fmt.Println("DelRobot5:", str23)

>>>>>>> origin/main
	_, str5 := script.JoinGame(1)
	fmt.Println("JoinGame1:", str5)
	_, str6 := script.JoinGame(2)
	fmt.Println("JoinGame2:", str6)
	_, str7 := script.JoinGame(0)
	fmt.Println("JoinGame0:", str7)

	_, str8 := script.LeaveGame()
	fmt.Println("LeaveGame:", str8)

	_, str9 := script.SetOrCancelPrepare()
	fmt.Println("SetOrCancelPrepare:", str9)
	_, str10 := script.GameStart()
	fmt.Println("GameStart:", str10)

	_, str12 := script.ChangePos(1)
	fmt.Println("ChangePos3:", str12)
	_, str13 := script.ChangePos(2)
	fmt.Println("ChangePos4:", str13)
	_, str14 := script.ChangePos(1)
	fmt.Println("ChangePos3:", str14)
	_, str15 := script.ChangePos(2)
	fmt.Println("ChangePos4:", str15)

<<<<<<< HEAD
	//_, str4 := script.DestroyRoom(roomId)
	//fmt.Println("DestroyRoom：", str4)
=======
	_, str16 := script.Act1GameInitBuf()
	fmt.Println("Act1GameInitBuf:", str16)

	_, str26 := script.TimeOut(0)
	fmt.Println("TimeOut0:", str26)
	_, str32 := script.TimeOut(1)
	fmt.Println("TimeOut1:", str32)
	_, str28 := script.TimeOut(2)
	fmt.Println("TimeOut2:", str28)
	_, str29 := script.TimeOut(3)
	fmt.Println("TimeOut3:", str29)
	_, str30 := script.TimeOut(4)
	fmt.Println("TimeOut4:", str30)
	_, str31 := script.TimeOut(5)
	fmt.Println("TimeOut5:", str31)

	_, str170 := script.PlayCard(0)
	fmt.Println("PlayCard0:", str170)
	_, str171 := script.PlayCard(1)
	fmt.Println("PlayCard1:", str171)
	_, str172 := script.PlayCard(2)
	fmt.Println("PlayCard2:", str172)
	_, str173 := script.PlayCard(3)
	fmt.Println("PlayCard3:", str173)
	_, str174 := script.PlayCard(4)
	fmt.Println("PlayCard4:", str174)
	_, str175 := script.PlayCard(5)
	fmt.Println("PlayCard5:", str175)

	_, str180 := script.EventHandle(0)
	fmt.Println("EventHandle0:", str180)
	_, str181 := script.EventHandle(1)
	fmt.Println("EventHandle1:", str181)
	_, str182 := script.EventHandle(2)
	fmt.Println("EventHandle2:", str182)
	_, str183 := script.EventHandle(3)
	fmt.Println("EventHandle3:", str183)
	_, str184 := script.EventHandle(4)
	fmt.Println("EventHandle4:", str184)
	_, str185 := script.EventHandle(5)
	fmt.Println("EventHandle5:", str185)

	_, str190 := script.GetCardFromPool()
	fmt.Println("GetCardFromPool:", str190)
>>>>>>> origin/main
}
