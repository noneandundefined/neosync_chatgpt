package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

func NewRedisDb() error {
	Client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password:        os.Getenv("REDIS_PASSWORD"),
		DB:              0,
		PoolSize:        35,
		MinIdleConns:    7,
		PoolTimeout:     17 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	fmt.Printf("[%v] [INFO] Successfully connected to redis\n", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func cacheKey(name string) string {
	return "cache:" + name
}
