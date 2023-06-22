package server

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"hdyx/common"
	"time"
)

type redisCfgObj struct {
	User         string
	Pass         string
	MaxOpenCon   int
	MaxIdleConns int
	ConLiveTime  int
	PieceNum     int
	AllHostCfg   map[string]any
}

var (
	RedisPools []*redis.Pool
	RedisCfg   redisCfgObj
)

func init() {
	// 配置初始化
	redisCfgIni()
	// redis初始化
	redisPoolInit()
	fmt.Println("--Redis初始化完成")
}

func redisCfgIni() {
	RedisCfg.User = common.GetYamlMapCfg("redisCfg", "redis", "common", "user").(string)
	RedisCfg.Pass = common.GetYamlMapCfg("redisCfg", "redis", "common", "pass").(string)
	RedisCfg.MaxOpenCon = common.GetYamlMapCfg("redisCfg", "redis", "common", "maxOpenCon").(int)
	RedisCfg.MaxIdleConns = common.GetYamlMapCfg("redisCfg", "redis", "common", "maxIdleConns").(int)
	RedisCfg.ConLiveTime = common.GetYamlMapCfg("redisCfg", "redis", "common", "conLiveTime").(int)
	RedisCfg.PieceNum = common.GetYamlMapCfg("redisCfg", "redis", "common", "pieceNum").(int)
	RedisCfg.AllHostCfg = common.GetYamlMapCfg("redisCfg", "redis", "game").(map[string]any)
}

func redisPoolInit() {
	for i := 1; i <= RedisCfg.PieceNum; i++ {
		host := fmt.Sprintf("%s%d", "host", i)
		var hostCfg map[string]any
		hostCfg = RedisCfg.AllHostCfg[host].(map[string]any)

		ip := hostCfg["ip"].(string)
		port := fmt.Sprintf("%d", hostCfg["port"])

		// 创建 Redis 连接池
		redisPool := &redis.Pool{
			MaxIdle:     RedisCfg.MaxIdleConns,                             // 最大空闲连接数
			MaxActive:   RedisCfg.MaxOpenCon,                               // 最大连接数
			IdleTimeout: time.Duration(RedisCfg.ConLiveTime) * time.Minute, // 空闲连接超时时间
			Dial: func() (redis.Conn, error) {
				// 建立连接的函数
				c, err := redis.Dial("tcp", ip+":"+port, redis.DialPassword(RedisCfg.Pass))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				// 连接测试函数
				_, err := c.Do("PING")
				return err
			},
		}

		RedisPools = append(RedisPools, redisPool)
	}
}
