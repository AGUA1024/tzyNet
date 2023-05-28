package model

import (
	"encoding/json"
	"hdyx/common"
	"strconv"
)

const (
	KEY_UIDTOROOM    = "uidToRoomId"
	KEY_ROOMIDTOINFO = "roomIdToRoomInfo"
)

type roomModel struct {
	actId      uint32          `json:"actId"`
	uidToState map[uint64]bool `json:"uidToState"`
}

func CreateRoom(ctx *common.ConContext, roomId uint64, actId uint32) bool {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	roomData := roomModel{
		actId: actId,
		uidToState: map[uint64]bool{
			ctx.GetConGlobalObj().Uid: false,
		},
	}

	data, _ := json.Marshal(roomData)

	ok := redis.RedisWrite(ctx, REDIS_ROOM, "HSET", KEY_ROOMIDTOINFO, strRoomId, string(data))

	return ok
}

func IsRoomExits(ctx *common.ConContext, roomId uint64) (bool, error) {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, err := redis.RedisQuery("HGET", KEY_ROOMIDTOINFO, strRoomId)
	if err != nil {
		return false, err
	}

	return data != nil, nil
}
