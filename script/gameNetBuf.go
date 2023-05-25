package script

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	ioBuf2 "hdyx/net/ioBuf"
)

func ConGlobalObjInitBuf() ([]byte, string) {
	buf := ConGlobalObjInit_InObj{
		Uid:    12,
		RoomId: 123,
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
	binaryString := byteSliceToBinaryString(data)

	return data, binaryString
}

func byteSliceToBinaryString(bytes []byte) string {
	var result string
	for _, b := range bytes {
		result += fmt.Sprintf("%08b", b)
	}
	return result
}
