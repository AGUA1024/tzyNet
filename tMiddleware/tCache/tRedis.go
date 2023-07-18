package tCache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"tzyNet/tIMiddleware"
)

type Redis struct {
	clusterClient *redis.ClusterClient
}

type RedisOpts struct {
	clusterHosts []string
}

func (this RedisOpts) SetClusterHosts(hosts []string) {
	this.clusterHosts = hosts
}

func (this RedisOpts) GetClusterHosts() []string {
	return this.clusterHosts
}

type RedisTxOperator struct {
	pipline redis.Pipeliner
}

func (this *Redis) NewCache(opts tIMiddleware.ICacheOpts) (tIMiddleware.ICache, error) {
	clusterCli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: opts.GetClusterHosts(),
	})

	return &Redis{clusterClient: clusterCli}, nil
}

func (this *Redis) Do(ctx context.Context, command string, args ...any) (any, error) {
	return this.clusterClient.Do(ctx, command, args).Result()
}

func (this *Redis) Multi() *RedisTxOperator {
	return &RedisTxOperator{
		pipline: this.clusterClient.TxPipeline(),
	}
}

func (this *RedisTxOperator) TxDo(ctx context.Context, command string, args ...any) (any, error) {
	// 执行Lua脚本
	script := `
		local result = redis.call(KEYS[1], unpack(ARGV))
		return result
	`

	this.pipline.Eval(ctx, script, []string{command}, args).Result()

	// 封装的通用函数执行Lua脚本
	return this.pipline.Eval(ctx, script, []string{command}, args).Result()
}

func (this *RedisTxOperator) TxExec() ([]redis.Cmder, error) {
	// 执行事务
	context.Background()
	return this.pipline.Exec(context.Background())
}
