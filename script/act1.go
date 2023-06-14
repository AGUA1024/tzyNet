package script

import (
	"encoding/hex"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	ioBuf2 "hdyx/net/ioBuf"
)

func Act1GameInitBuf() ([]byte, string) {
	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30001,
		Data:           nil,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	// 将 byte 装换为 2进制的字符串
	binaryString := hex.EncodeToString(data)

	return data, binaryString
}

func PlayCard(cardIndex uint32) ([]byte, string) {
	buf := PlayCard_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		CardIndex:     cardIndex,
	}

	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30002,
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

func EventHandle(chooseIndex uint32) ([]byte, string) {
	buf := EventHandle_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		ChooseIndex:   chooseIndex,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30003,
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

func GetCardFromPool() ([]byte, string) {
	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30004,
		Data:           nil,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	// 将 byte 装换为 2进制的字符串
	binaryString := hex.EncodeToString(data)

	return data, binaryString
}

func TimeOut(seqId uint32) ([]byte, string) {
	buf := TimeOut_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		SeqId:         seqId,
	}
	cendBuf, _ := proto.Marshal(&buf)
	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30005,
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

func PlayerReCon() ([]byte, string) {
	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x30006,
		Data:           nil,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	// 将 byte 装换为 2进制的字符串
	binaryString := hex.EncodeToString(data)

	return data, binaryString
}
