package tIMiddleware

import (
	"context"
	"tzyNet/tMiddleware"
)

type ICache interface {
	NewCache(opts ICacheOpts) (ICache, error)
	Do(ctx context.Context, command string, args ...any) (any, error)
}

type ICacheOpts interface {
	tMiddleware.OptsBase
}
