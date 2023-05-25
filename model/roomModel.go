package model

import (
	"encoding/json"
	"hdyx/common"
	"hdyx/global"
	"hdyx/server"
)

const (
	KEY_UIDTOROOM    = "uidToRoomId"
	KEY_ROOMIDTOINFO = "roomIdToRoomInfo"
)

type roomMaster struct {
	mpUidToRoomId    map[uint64]uint64
	mpRoomIdToActId  map[uint64]uint32
	roomIdToRoomInfo map[uint64]map[uint64]bool
}

var RoomMaster = roomMaster{
	mpUidToRoomId:    map[uint64]uint64{},
	mpRoomIdToActId:  map[uint64]uint32{},
	roomIdToRoomInfo: map[uint64]map[uint64]bool{},
}

func (m *roomMaster) CreateRoom(uid uint64, roomId uint64, actId uint32) {
	m.mpUidToRoomId[uid] = roomId
	m.mpRoomIdToActId[roomId] = actId
	m.roomIdToRoomInfo[roomId] = map[uint64]bool{uid: false}
}

func (m *roomMaster) RegisterUidToRoom(uid uint64, roomId uint64, actId uint32) {
	m.mpUidToRoomId[uid] = roomId

	if _, ok := m.roomIdToRoomInfo[roomId]; !ok {
		// 载入数据
		m.loadRoomInfo(roomId)
	}
}

func (m *roomMaster) IsRoomExits(ctx *global.ConContext, roomId uint64) bool {
	redis := server.GetRedis(roomId)
	jsData, err := redis.RedisDo("HGET", KEY_ROOMIDTOINFO, roomId)
	if err != nil {
		common.Logger.GameErrorLog(ctx, common.ERR_REDIS_LOAD_ROOMINFO, "redis获取房间id失败")
	}

	return jsData != nil
}

func (m *roomMaster) loadRoomInfo(roomId uint64) {
	redis := server.GetRedis(roomId)
	jsData, err := redis.RedisDo("HGET", KEY_ROOMIDTOINFO, roomId)
	if err != nil {
		common.Logger.SystemErrorLog("LOAD_ROOMINFO_ERROR", err)
	}
	data, _ := jsData.([]byte)
	err = json.Unmarshal(data, &m.roomIdToRoomInfo)
}
