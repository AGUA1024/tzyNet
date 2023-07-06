package tCache

import (
	"github.com/gomodule/redigo/redis"
	"tzyNet/tIMiddleware"
)

type tRedis struct {
	conPool *redis.Pool
}

func newRedis(address string, userName string, password string) (tIMiddleware.ICache, error) {

	redisPool := &redis.Pool{
		Dial:            nil,
		DialContext:     nil,
		TestOnBorrow:    nil,
		MaxIdle:         0,
		MaxActive:       0,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	redis.

	con, err := redis.Dial("tcp", address, redis.DialUsername(userName), redis.DialPassword(password))
	if err != nil {
		return nil, err
	}

	return &tRedis{
		con: &con,
	}, nil
}

func (this *tRedis) Do(command string, args ...any) (any, error) {
	return nil, nil
}
