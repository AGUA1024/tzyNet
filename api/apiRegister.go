package api

import "google.golang.org/protobuf/proto"

var masterRouteRegister = map[uint32]map[uint32]func([]byte){
	0x1: userRouteRegMp,
}

var userRouteRegMp = map[uint32]func([]byte){
	0x1: GetUserInfo,
}

func GetApiByCmd(cmd uint32) func([]byte) {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff

	return masterRouteRegister[highCmd][lowCmd]
}

func getParamObj[T proto.Message](params []byte, obj T) T {
	proto.Unmarshal(params, obj)
	return obj
}
