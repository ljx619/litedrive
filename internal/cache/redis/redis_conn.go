package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"litedrive/internal/utils"
	"log"
	"time"
)

var (
	Ctx      = context.Background()
	RedisCli *redis.Client
)

// InitRedis 初始化 Redis 连接
func InitRedis() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	RedisCli = redis.NewClient(&redis.Options{
		Addr:         config.Redis.Host,     // Redis 地址
		Password:     config.Redis.Password, // 没有密码则为空
		DB:           0,                     // 默认数据库
		PoolSize:     20,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
	})

	// 测试连接
	if _, err = RedisCli.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}

// GetClient 获取 Redis 客户端，确保 client 被初始化
func GetClient() (*redis.Client, error) {
	if RedisCli == nil {
		return nil, errors.New("Redis client is not initialized")
	}
	return RedisCli, nil
}

// Set 设置 key-value，并可指定过期时间
func Set(key string, value string, expiration time.Duration) error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	return c.Set(Ctx, key, value, expiration).Err()
}

// Get 获取 key 的值
func Get(key string) (string, error) {
	c, err := GetClient()
	if err != nil {
		return "", err
	}
	return c.Get(Ctx, key).Result()
}

// Del 删除 key
func Del(key string) error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	return c.Del(Ctx, key).Err()
}

// Exists 检查 key 是否存在
func Exists(key string) (bool, error) {
	c, err := GetClient()
	if err != nil {
		return false, err
	}
	n, err := c.Exists(Ctx, key).Result()
	return n > 0, err
}

// Close 关闭 Redis 连接
func CloseRedis() {
	if RedisCli == nil {
		log.Println("Redis client is not initialized")
		return
	}
	_ = RedisCli.Close()
	log.Println("Redis closed")
}
