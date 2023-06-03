package model

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	api "hdyx/api/protobuf"
	"hdyx/common"
	"hdyx/net/ioBuf"
	"hdyx/sdk"
	"strconv"
)

const (
	KEY_UIDTOROOM    = "uidToRoomId"
	KEY_ROOMIDTOINFO = "roomIdToRoomInfo"
)

type roomModel struct {
	ActId          uint32                  `json:"ActId"`
	GameLv         uint32                  `json:"GameLv"`
	PosIdToPlayer  map[uint32]*PlayerModel `json:"PosIdToPlayer"`
	ArrUidAudience []uint64                `json:"ArrUidAudience"`
}

type PlayerModel struct {
	Uid      uint64
	IsMaster bool // 房主
	State    bool // 是否准备
	Head     string
	Name     string
}

func getRoomIdToInfoKey(ctx *common.ConContext) string {
	return GetRedisPreKey(ctx.GetConGlobalObj().RoomId) + KEY_ROOMIDTOINFO
}

func CreateRoom(ctx *common.ConContext, user *sdk.UserInfo, roomId uint64, actId uint32, gameLv uint32) *api.CreateRoom_OutObj {
	redis := GetCacheById(roomId)
	strRoomId := strconv.FormatUint(roomId, 10)

	uid := ctx.GetConGlobalObj().Uid

	roomData := roomModel{
		ActId:  actId,
		GameLv: gameLv,
		PosIdToPlayer: map[uint32]*PlayerModel{
			0: {
				Uid:      uid,
				IsMaster: true,
				State:    false,
				Head:     user.Cover,
				Name:     user.UserName,
			},
		},
		ArrUidAudience: []uint64{},
	}

	data, _ := json.Marshal(roomData)

	ok := redis.RedisWrite(ctx, REDIS_ROOM, "HSET", getRoomIdToInfoKey(ctx), strRoomId, string(data))
	if !ok {
		return nil
	}

	ret := &api.CreateRoom_OutObj{
		Players: map[uint32]*api.PlayerInfo{
			0: {
				Uid:      uid,
				IsMaster: true,
				State:    false,
				Head:     user.Cover,
				Name:     user.UserName,
			},
		},
	}

	return ret
}

func LeaveRoomAndBroadcast(ctx *common.ConContext) bool {
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomInfo, err := GetGameRoomInfo(ctx, roomId)
	if roomInfo == nil || err != nil {
		return false
	}

	// 是否在观众席中
	ok, roomIndex := GetRoomIndex(ctx, roomInfo)
	if ok {
		// 离开观众席位
		roomInfo.ArrUidAudience = append(roomInfo.ArrUidAudience[:roomIndex], roomInfo.ArrUidAudience[roomIndex+1:]...)
	} else {
		playerInfo := roomInfo.PosIdToPlayer
		// 是否在房间
		ok, gameIndex := GetGameIndex(ctx, roomInfo)
		if !ok {
			return true
		}

		if isMaster := playerInfo[gameIndex].IsMaster; isMaster {
			// 房主转移
			for i, player := range playerInfo {
				if i == gameIndex {
					continue
				}

				// 设置新房主
				if player.IsMaster == false {
					playerInfo[i].IsMaster = true
					break
				}
			}
		}

		// 离开游戏房间
		delete(playerInfo, gameIndex)
		roomInfo.PosIdToPlayer = playerInfo
	}

	// 离开后如果没有玩家，则销毁房间
	if len(roomInfo.PosIdToPlayer) == 0 && len(roomInfo.ArrUidAudience) == 0 {
		DestroyRoom(ctx, roomId)
		return true
	}

	// 更新房间数据写入redis
	RoomModelSave(ctx, roomInfo)

	// 广播
	var mpRet = BackGameRoomInfo(roomInfo.PosIdToPlayer)
	broadCastInfo := &api.LeaveGame_OutObj{Players: mpRet}
	MsgRoomBroadcast[*api.LeaveGame_OutObj](ctx, broadCastInfo)

	return true
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

func GetRoomIndex(ctx *common.ConContext, GameRoomInfo *roomModel) (bool, int) {
	for index, uid := range GameRoomInfo.ArrUidAudience {
		if uid == ctx.GetConGlobalObj().Uid {
			return true, index
		}
	}

	return false, 9999999
}

func GetGameIndex(ctx *common.ConContext, GameRoomInfo *roomModel) (bool, uint32) {
	for index, player := range GameRoomInfo.PosIdToPlayer {
		if player.Uid == ctx.GetConGlobalObj().Uid {
			return true, index
		}
	}

	return false, 999
}

func GetGameRoomInfo(ctx *common.ConContext, roomId uint64) (*roomModel, error) {
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
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)

	// 玩家广播
	for _, player := range roomInfo.PosIdToPlayer {
		uid := player.Uid
		_, ok := common.MpUserStorage[uid]
		if !ok {
			continue
		}

		if err = common.MpUserStorage[uid].WsCon.WriteMessage(common.BinaryMessage, outStream); err != nil {
			continue
		}
	}

	// 观众广播
	for _, uid := range roomInfo.ArrUidAudience {
		_, ok := common.MpUserStorage[uid]
		if !ok {
			continue
		}

		if err = common.MpUserStorage[uid].WsCon.WriteMessage(common.BinaryMessage, outStream); err != nil {
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

func BackGameRoomInfo(roomInfo map[uint32]*PlayerModel) map[uint32]*api.PlayerInfo {
	var mpRet = map[uint32]*api.PlayerInfo{}

	for index, playerInfo := range roomInfo {
		mpRet[index] = &api.PlayerInfo{
			Uid:      playerInfo.Uid,
			State:    playerInfo.State,
			IsMaster: playerInfo.IsMaster,
			Head:     playerInfo.Head,
			Name:     playerInfo.Name,
		}
	}

	return mpRet
}
