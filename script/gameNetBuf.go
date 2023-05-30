package script

import (
	"encoding/hex"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	ioBuf2 "hdyx/net/ioBuf"
)

func GetGateWayBuf() ([]byte, string) {
	buf := GetGateWay_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        11111,
	}
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

	// 将 byte 装换为 2进制的字符串
	binaryString := hex.EncodeToString(data)

	return data, binaryString
}

func ConGlobalObjInitBuf(uid uint64, roomId uint64) ([]byte, string) {
	buf := ConGlobalObjInit_InObj{
		Uid:    uid,
		RoomId: roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x10002,
		Data:           cendBuf,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	// 将 byte 装换为 2进制的字符串
	binaryString := hex.EncodeToString(data)

	return data, binaryString
}

//func byteSliceToBinaryString(bytes []byte) string {
//	var result string
//	for _, b := range bytes {
//		result += fmt.Sprintf("%08b", b)
//	}
//	return result
//}
