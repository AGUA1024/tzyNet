package server

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"hdyx/common"
	"log"
	"time"
)

var redisPools []*redis.Pool

type redisCfgObj struct {
	User         string
	Pass         string
	MaxOpenCon   int
	MaxIdleConns int
	ConLiveTime  int
	PieceNum     int
	AllHostCfg   map[string]any
}

type RedisOperater struct {
	pool redis.Pool
}

func GetRedis(uid int32) RedisOperater {
	piece := uid % int32(redisCfg.PieceNum)
	var redisOp = RedisOperater{
		pool: *redisPools[piece],
	}

	return redisOp
}

var redisCfg redisCfgObj

func init() {
	//// 配置初始化
	//redisCfgIni()
	//// redis初始化
	//redisPoolInit()
}

func redisCfgIni() {
	redisCfg.User = common.GetYamlMapCfg("redisCfg", "redis", "common", "user").(string)
	redisCfg.Pass = common.GetYamlMapCfg("redisCfg", "redis", "common", "pass").(string)
	redisCfg.MaxOpenCon = common.GetYamlMapCfg("redisCfg", "redis", "common", "maxOpenCon").(int)
	redisCfg.MaxIdleConns = common.GetYamlMapCfg("redisCfg", "redis", "common", "maxIdleConns").(int)
	redisCfg.ConLiveTime = common.GetYamlMapCfg("redisCfg", "redis", "common", "conLiveTime").(int)
	redisCfg.PieceNum = common.GetYamlMapCfg("redisCfg", "redis", "common", "pieceNum").(int)
	redisCfg.AllHostCfg = common.GetYamlMapCfg("redisCfg", "redis", "game").(map[string]any)
}

func redisPoolInit() {
	for i := 1; i <= redisCfg.PieceNum; i++ {
		host := fmt.Sprintf("%s%d", "host", i)
		var hostCfg map[string]any
		hostCfg = dbCfg.AllHostCfg[host].(map[string]any)

		ip := hostCfg["ip"].(string)
		port := fmt.Sprintf("%d", hostCfg["port"])

		// 创建 Redis 连接池
		redisPool := &redis.Pool{
			MaxIdle:     redisCfg.MaxIdleConns,                             // 最大空闲连接数
			MaxActive:   redisCfg.MaxOpenCon,                               // 最大连接数
			IdleTimeout: time.Duration(redisCfg.ConLiveTime) * time.Minute, // 空闲连接超时时间
			Dial: func() (redis.Conn, error) {
				// 建立连接的函数
				c, err := redis.Dial("tcp", ip+":"+port)
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

		redisPools = append(redisPools, redisPool)
	}
}

// 查询数据
func (op *RedisOperater) redisDo(commandName string, args ...interface{}) (any, error) {
	// 使用连接池获取连接并进行操作
	conn := op.pool.Get()
	defer conn.Close()

	data, err := conn.Do(commandName, args)
	if err != nil {
		log.Fatal(err)
	}

	return data, err
}
