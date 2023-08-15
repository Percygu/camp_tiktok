package lock

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	pool "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redis "github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"sync"
	"time"
	"usersvr/config"
)

var (
	rs          *redsync.Redsync
	lockClients []*redis.Client
	mutexOnce   sync.Once
	lockExpiry  = time.Second * 10
	retryDelay  = time.Millisecond * 100
	tries       = 3
	// 设置锁的过期时间
	// 设置重试次数
	// 设置重试间隔
	option = []redsync.Option{
		redsync.WithExpiry(lockExpiry),
		redsync.WithRetryDelay(retryDelay),
		redsync.WithTries(tries),
	}
	lockPrefix = "tiktok:lock:"
)

// openDB 连接db
func initLock() {
	// 初始化多台 redis master连接
	var pools []pool.Pool
	for _, conf := range config.GetGlobalConfig().RedsyncConfig {
		lockClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Password: conf.Password,
			PoolSize: conf.PoolSize,
		})

		if _, err := lockClient.Ping(context.Background()).Result(); err != nil {
			panic("Failed to ping redisMutex, err:%s")
		}
		// 放入
		lockClients = append(lockClients, lockClient)
		pools = append(pools, goredis.NewPool(lockClient))
	}

	rs = redsync.New(pools...)
}

func CloseLock() {
	for _, client := range lockClients {
		client.Close()
	}
}

// 获取锁
func GetLock(name string) *redsync.Mutex {
	// 初始化一次
	mutexOnce.Do(initLock)
	return rs.NewMutex(lockPrefix+name, option...)
}

func UnLock(mutex *redsync.Mutex) {
	mutex.Unlock()
}
