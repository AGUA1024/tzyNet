package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"tzyNet/tCommon"
	"tzyNet/tServer"
)

type RedisEventQueue struct {
	events []*RedisEvent
}

type RedisEvent struct {
	command string
	args    []interface{}
}

type RedisOperator struct {
	pool *redis.Pool
}

const (
	REDIS_PLAYER = 1
	REDIS_ROOM   = 2
)

func GetRedisPreKey(id uint64) string {
	return fmt.Sprintf("%02d_", id%100)
}

func (this RedisEvent) GetCommand() string {
	return this.command
}

func (this RedisEvent) GetArgs() []interface{} {
	return this.args
}

func AllRedisSave(ctx *tCommon.ConContext) {
	ctx.AllCacheSave(GetCacheById(ctx.GetConGlobalObj().Uid), GetCacheById(ctx.GetConGlobalObj().RoomId))
}

// 内存数据落地redis
func (this *RedisOperator) CacheSave(arrCacheEvent []tCommon.CacheEvent) error {
	// 使用连接池获取连接并进行操作
	conn := this.pool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, event := range arrCacheEvent {
		err = conn.Send(event.GetCommand(), event.GetArgs()...)
		if err != nil {
			return err
		}
	}

	conn.Do("EXEC")
	fmt.Println("cache save")
	return err
}

func GetCacheById(id uint64) *RedisOperator {
	piece := id % uint64(tServer.RedisCfg.PieceNum)
	var redisOp = RedisOperator{
		pool: tServer.RedisPools[piece],
	}
	return &redisOp
}

func (this *RedisOperator) RedisQuery(command string, args ...interface{}) (any, error) {
	conn := this.pool.Get()
	defer conn.Close()

	data, err := conn.Do(command, args...)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (this *RedisOperator) RedisWrite(ctx *tCommon.ConContext, typeId int, command string, args ...interface{}) bool {
	var ok bool

	event := RedisEvent{
		command: command,
		args:    args,
	}

	switch typeId {
	case REDIS_PLAYER:
		ok = ctx.PlayerRedisEventPush(event)
	case REDIS_ROOM:
		ok = ctx.RoomRedisEventPush(event)
	}

	return ok
}
