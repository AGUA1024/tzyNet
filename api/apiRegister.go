package api

import "hdyx/common"

var masterRouteRegister = map[uint32]map[uint32]func(*common.ConContext, []byte){
	0x1: gameNetRegMp,
	0x2: roomRegMp,
}

var gameNetRegMp = map[uint32]func(*common.ConContext, []byte){
	0x1: GetGateWay,
	0x2: ConGlobalObjInit,
}

var roomRegMp = map[uint32]func(*common.ConContext, []byte){
	0x1: CreateRoom,
	0x2: DestroyRoom,
	0x3: JoinRoom,
	0x4: LeaveRoom,
	0x5: SetOrCancelPrepare,
	0x6: GameStart,
}

func GetApiByCmd(cmd uint32) func(*common.ConContext, []byte) {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff

	return masterRouteRegister[highCmd][lowCmd]
}
