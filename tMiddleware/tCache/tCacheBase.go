package tCache

import (
	"errors"
	"tzyNet/tIMiddleware"
)

const (
	CacheType_Redis uint16 = iota
)

func NewCache(cacheType uint16, address string, userName string, password string) (tIMiddleware.ICache, error) {
	switch cacheType {
	case CacheType_Redis:
		redis, err := newRedis(address, userName, password)
		return redis, err
	default:
		return nil, errors.New("Invalid cache")
	}
}
