package model

import (
	"encoding/json"
	"fmt"
	"hdyx/common"
)

type GameCfg struct {
	ActId        uint32
	MaxPlayerNum int
	MinPlayerNum int
}

type ActBaseInterface interface {
	NewActModel(ctx *common.ConContext) ActBaseInterface
	Save(ctx *common.ConContext) (bool, error)
	LoadModelInfo() (bool, error)
	GetActCfg() *GameCfg
}

type actBaseModel struct {
	RoomId  uint64
	ActId   uint32
	ActInfo map[string]any
	IsOver  bool
}

func (this *actBaseModel) GetActKey() string {
	return GetRedisPreKey(this.RoomId) + "act_" + fmt.Sprintf("%d", this.ActId)
}

func (this *actBaseModel) LoadModelInfo() (bool, error) {
	key := this.GetActKey()
	cache := GetCacheById(this.RoomId)
	data, err := cache.RedisQuery("HGET", key, this.RoomId)
	if data == nil || err != nil {
		return false, err
	}

	json.Unmarshal(data.([]byte), this)

	return true, nil
}

func (this *actBaseModel) Save(ctx *common.ConContext) (bool, error) {
	key := this.GetActKey()
	cache := GetCacheById(this.RoomId)

	arrByte, err := json.Marshal(this)
	if err != nil {
		return false, err
	}

	ok := cache.RedisWrite(ctx, REDIS_ROOM, "HSET", key, this.RoomId, string(arrByte))
	return ok, nil
}

func (this *actBaseModel) Destory(ctx *common.ConContext) bool {
	key := this.GetActKey()
	cache := GetCacheById(this.RoomId)

	ok := cache.RedisWrite(ctx, REDIS_ROOM, "HDEL", key, this.RoomId)
	return ok
}
