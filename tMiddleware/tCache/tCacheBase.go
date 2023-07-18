package tCache

import (
	"tzyNet/tIMiddleware"
)

const (
	CacheType_Redis uint16 = iota
)

func NewCache[cacheType tIMiddleware.ICache](opts tIMiddleware.ICacheOpts) (tIMiddleware.ICache, error) {
	var cache cacheType
	return cache.NewCache(opts)
}
