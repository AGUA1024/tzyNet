package model

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"hdyx/common"
	"hdyx/net/ioBuf"
	"strconv"
)

const (
	KEY_UIDTOROOM    = "uidToRoomId"
	KEY_ROOMIDTOINFO = "roomIdToRoomInfo"
)

type roomModel struct {
	ActId           uint32                  `json:"ActId"`
	UidToPlayerInfo map[uint64]*PlayerModel `json:"uidToState"`
}

type PlayerModel struct {
	IsMaster bool // 房主
	State    bool // 是否准备
}

func getRoomIdToInfoKey(ctx *common.ConContext) string {
	return GetRedisPreKey(ctx.GetConGlobalObj().RoomId) + KEY_ROOMIDTOINFO
}

func CreateRoom(ctx *common.ConContext, roomId uint64, actId uint32) bool {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	roomData := roomModel{
		ActId: actId,
		UidToPlayerInfo: map[uint64]*PlayerModel{
			ctx.GetConGlobalObj().Uid: {
				IsMaster: true,
				State:    false,
			},
		},
	}

	data, _ := json.Marshal(roomData)

	ok := redis.RedisWrite(ctx, REDIS_ROOM, "HSET", getRoomIdToInfoKey(ctx), strRoomId, string(data))

	return ok
}

func IsRoomExits(ctx *common.ConContext, roomId uint64) (bool, error) {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, err := redis.RedisQuery("HGET", getRoomIdToInfoKey(ctx), strRoomId)
	if err != nil {
		return false, err
	}

	return data != nil, nil
}

func IsRoomMaster(ctx *common.ConContext, roomId uint64) (bool, error) {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, err := redis.RedisQuery("HGET", getRoomIdToInfoKey(ctx), strRoomId)
	if err != nil || data == nil {
		return false, err
	}

	var roomInfo roomModel
	err = json.Unmarshal(data.([]byte), &roomInfo)
	if err != nil {
		return false, err
	}

	ok := roomInfo.UidToPlayerInfo[ctx.GetConGlobalObj().Uid].IsMaster
	return ok, nil
}

func IsInRoom(ctx *common.ConContext, roomId uint64) (bool, error) {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, err := redis.RedisQuery("HGET", getRoomIdToInfoKey(ctx), strRoomId)
	if err != nil || data == nil {
		return false, err
	}

	var roomInfo roomModel
	err = json.Unmarshal(data.([]byte), &roomInfo)
	if err != nil {
		return false, err
	}

	_, ok := roomInfo.UidToPlayerInfo[ctx.GetConGlobalObj().Uid]
	return ok, nil
}

func GetRoomInfo(ctx *common.ConContext, roomId uint64) (*roomModel, error) {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, err := redis.RedisQuery("HGET", getRoomIdToInfoKey(ctx), strRoomId)
	if err != nil || data == nil {
		return nil, err
	}

	var roomInfo roomModel
	json.Unmarshal(data.([]byte), &roomInfo)

	return &roomInfo, nil
}

func DestroyRoom(ctx *common.ConContext, roomId uint64) bool {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	ok := redis.RedisWrite(ctx, REDIS_ROOM, "HDEL", getRoomIdToInfoKey(ctx), strRoomId)

	return ok
}

func MsgRoomBroadcast[T proto.Message](ctx *common.ConContext, obj T) (any, error) {
	// 封装输出数据
	probuf, _ := proto.Marshal(obj)

	out := ioBuf.OutPutBuf{
		Uid:            ctx.GetConGlobalObj().Uid,
		CmdCode:        0,
		ProtocolSwitch: 0,
		CmdMerge:       ctx.GetConGlobalObj().Cmd,
		ResponseStatus: 0,
		Data:           probuf,
	}

	outStream, err := proto.Marshal(&out)
	if err != nil {
		common.Logger.SystemErrorLog("GET_OUT_STREAM_Marshal_ERROR", err)
	}

	// 获取广播列表
	roomInfo, err := GetRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)

	// 广播
	for uid, _ := range roomInfo.UidToPlayerInfo {
		_, ok := common.MpUserStorage[uid]
		if !ok {
			continue
		}

		if err = common.MpUserStorage[uid].WsCon.WriteMessage(common.TextMessage, outStream); err != nil {
			continue
		}
	}

	return true, nil
}

func RoomModelSave(ctx *common.ConContext, model *roomModel) bool {
	roomId := ctx.GetConGlobalObj().RoomId
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	data, _ := json.Marshal(model)

	ok := redis.RedisWrite(ctx, REDIS_ROOM, "HSET", getRoomIdToInfoKey(ctx), strRoomId, string(data))

	return ok
}
