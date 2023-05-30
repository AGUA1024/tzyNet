package script

import (
	"encoding/hex"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	ioBuf2 "hdyx/net/ioBuf"
)

func CreateRoomBuf(roomId uint64) ([]byte, string) {
	buf := CreateRoom_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
		ActId:         1,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20001,
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

func DestroyRoom(roomId uint64) ([]byte, string) {
	buf := DestroyRoom_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20002,
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

func JoinRoom(roomId uint64) ([]byte, string) {
	buf := JoinRoom_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20003,
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

func LeaveRoom(roomId uint64) ([]byte, string) {
	buf := LeaveRoom_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20004,
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

func SetOrCancelPrepare(roomId uint64) ([]byte, string) {
	buf := SetOrCancelPrepare_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20005,
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

func GameStart(roomId uint64) ([]byte, string) {
	buf := GameStart_InObj{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RoomId:        roomId,
	}
	cendBuf, _ := proto.Marshal(&buf)

	// 发送消息
	message := &ioBuf2.ClientBuf{
		ProtocolSwitch: 1,
		CmdMerge:       0x20006,
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
