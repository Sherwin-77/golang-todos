package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sherwin-77/go-echo-template/configs"
)

func InitRedis(config configs.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}

type Cache interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) string
	Del(key string) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{client}
}

func (c *cache) Set(key string, value interface{}, duration time.Duration) error {
	return c.client.Set(context.Background(), key, value, duration).Err()
}

func (c *cache) Get(key string) string {
	value, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}

	return value
}

func (c *cache) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}
