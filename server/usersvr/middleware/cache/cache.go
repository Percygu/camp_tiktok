package cache

import (
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"sync"
	"time"
	"usersvr/config"
	"usersvr/log"
)

var (
	redisConn   *redis.Client
	redisOnce   sync.Once
	ValueExpire = time.Hour * 24 * 7
)

// openDB 连接db
func initRedis() {
	redisConfig := config.GetGlobalConfig().RedisConfig
	log.Infof("redisConfig=======%+v", redisConfig)
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PassWord,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	res, err := redisConn.Set(context.Background(), "abc", 100, 60).Result()
	log.Infof("res=======%v,err======%v", res, err)
	_, err = redisConn.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping redis, err:%s")
	}
}

func CloseRedis() {
	redisConn.Close()
}

// GetRedisCli 获取数据库连接
func GetRedisCli() *redis.Client {
	redisOnce.Do(initRedis)

	return redisConn
}

func CacheGetUser(key string) (string, error) {
	return "", nil
}
