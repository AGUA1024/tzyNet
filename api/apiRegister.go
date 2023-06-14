package api

import (
	"hdyx/common"
)

var masterRouteRegister = map[uint32]map[uint32]func(*common.ConContext, []byte){
	0x1: gameNetRegMp,
	0x2: roomRegMp,
	0x3: act1RegMp,
}

var gameNetRegMp = map[uint32]func(*common.ConContext, []byte){
	0x1: GetGateWay,
	0x2: ConGlobalObjInit,
}

var roomRegMp = map[uint32]func(*common.ConContext, []byte){
	0x1:  IsRoomExist,
	0x2:  CreateRoom,
	0x3:  JoinRoom,
	0x4:  JoinGame,
	0x5:  ChangePos,
	0x6:  LeaveGame,
	0x7:  SetOrCancelPrepare,
	0x8:  GameStart,
	0x9:  AddRobot,
	0x10: DelRobot,
	0x11: GetRoomInfo,
}

var act1RegMp = map[uint32]func(*common.ConContext, []byte){
	0x1: Act1GameInit,
	0x2: PlayCard,
	0x3: EventHandle,
	0x4: GetCardFromPool,
	0x5: TurnTimeOut,
	0x6: GetAct1Info,
}

func GetApiByCmd(cmd uint32) func(*common.ConContext, []byte) {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff

	return masterRouteRegister[highCmd][lowCmd]
}
