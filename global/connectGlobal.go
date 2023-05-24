package global

import "github.com/gorilla/websocket"

// connectGorutine局部空间
var MpConRoutineStorage = map[uint64]RoutineStorage{}

type RoutineStorage struct {
	WsCon  *websocket.Conn
	Uid    uint64
	Cmd    uint32
	RoomId uint64
}

// connectGorutine上下文
type ConContext struct {
	ConnectId uint64
}

func (this ConContext) GetConnectId() uint64 {
	return this.ConnectId
}

// readOnly
func (ctx ConContext) GetConGlobalVal() *RoutineStorage {
	conGlobalObj, ok := MpConRoutineStorage[ctx.GetConnectId()]
	if !ok {
		return nil
	}
	return &conGlobalObj
}

func (ctx ConContext) SetConGlobalWsCon(conn *websocket.Conn) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.WsCon = conn

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx ConContext) SetConGlobalCmd(cmd uint32) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.Cmd = cmd

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx ConContext) SetConGlobalRoomId(roomId uint64) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}

	conGlobalObj.RoomId = roomId

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx ConContext) SetConGlobalUid(uid uint64) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.Uid = uid

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}
