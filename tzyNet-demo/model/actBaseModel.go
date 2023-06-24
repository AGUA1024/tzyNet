package model

import (
	"encoding/json"
	"fmt"
	"tzyNet/tCommon"
)

type GameCfg struct {
	ActId        uint32
	MaxPlayerNum int
	MinPlayerNum int
}

type ActBaseModel struct {
	RoomId uint64
	ActId  uint32
	IsOver bool
}

type ActModelInterface interface {
	GetActId() uint32
	GetRoomId() uint64
	GetActCfg() *GameCfg
	IsPlayer(ctx *tCommon.ConContext) bool
	PlayerReConn(ctx *tCommon.ConContext) bool
}

func GetActKey(actId uint32, roomId uint64) string {
	return GetRedisPreKey(roomId) + "act_" + fmt.Sprintf("%d", actId)
}

// 加载游戏信息
func GetActModel(ctx *tCommon.ConContext, actId uint32) ActModelInterface {
	actModel := actModelRegister[actId]
	ok := LoadModelInfo(ctx, actId, actModel)
	if !ok {
		return nil
	}
	fmt.Println("act:", actModel)
	return actModel
}

func LoadModelInfo(ctx *tCommon.ConContext, actId uint32, actModel ActModelInterface) bool {
	roomId := ctx.GetConGlobalObj().RoomId
	key := GetActKey(actId, roomId)
	cache := GetCacheById(roomId)

	fmt.Println("HGET", key, roomId)

	data, err := cache.RedisQuery("HGET", key, roomId)
	if data == nil || err != nil {
		return false
	}

	json.Unmarshal(data.([]byte), actModel)
	fmt.Println("this:", actModel)
	return true
}

func Save(ctx *tCommon.ConContext, actModel ActModelInterface) (bool, error) {
	modelActId := actModel.GetActId()
	roomId := actModel.GetRoomId()

	key := GetActKey(modelActId, roomId)
	cache := GetCacheById(roomId)

	arrByte, err := json.Marshal(actModel)
	if err != nil {
		return false, err
	}

	ok := cache.RedisWrite(ctx, REDIS_ROOM, "HSET", key, roomId, string(arrByte))
	return ok, nil
}

func Destory(ctx *tCommon.ConContext, actModel ActModelInterface) bool {
	modelActId := actModel.GetActId()
	roomId := actModel.GetRoomId()

	key := GetActKey(modelActId, roomId)
	cache := GetCacheById(roomId)

	ok := cache.RedisWrite(ctx, REDIS_ROOM, "HDEL", key, roomId)
	return ok
}

func (this *ActBaseModel) GetRoomId() uint64 {
	return this.RoomId
}

func (this *ActBaseModel) GetActId() uint32 {
	return this.ActId
}
