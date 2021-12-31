package cache

import (
	"context"
	"errors"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/multi"
)

// Init 初始化缓存服务
func Init() error {
	universalOptions := &redis.UniversalOptions{
		Addrs:       strings.Split(serverConsts.CONFIG.Redis.Addr, ","),
		Password:    serverConsts.CONFIG.Redis.Password,
		PoolSize:    serverConsts.CONFIG.Redis.PoolSize,
		IdleTimeout: 300 * time.Second,
	}
	serverConsts.CACHE = redis.NewUniversalClient(universalOptions)
	err := multi.InitDriver(
		&multi.Config{
			DriverType:      serverConsts.CONFIG.System.CacheType,
			UniversalClient: serverConsts.CACHE},
	)
	if err != nil {
		return err
	}
	if multi.AuthDriver == nil {
		return errors.New("初始化认证驱动失败")
	}

	return nil
}

// SetCache 缓存数据
func SetCache(key string, value interface{}, expiration time.Duration) error {
	err := serverConsts.CACHE.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteCache 删除缓存数据
func DeleteCache(key string) (int64, error) {
	return serverConsts.CACHE.Del(context.Background(), key).Result()
}

// GetCacheString 获取字符串类型数据
func GetCacheString(key string) (string, error) {
	return serverConsts.CACHE.Get(context.Background(), key).Result()
}

// GetCacheBytes 获取bytes类型数据
func GetCacheBytes(key string) ([]byte, error) {
	return serverConsts.CACHE.Get(context.Background(), key).Bytes()
}

// GetCacheUint 获取uint类型数据
func GetCacheUint(key string) (uint64, error) {
	return serverConsts.CACHE.Get(context.Background(), key).Uint64()
}
