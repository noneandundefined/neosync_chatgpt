package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func configurationErrKey(imei string) string {
	return cacheKey("cfg-error:imei:" + imei)
}

func configurationKey(imei string) string {
	return cacheKey("cfg:imei:" + imei)
}

func SetCacheConfiguration(imei string, ttl time.Duration, cfg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Set(ctx, configurationKey(imei), cfg, ttl).Err()
}

func GetCacheConfiguration(imei string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := Client.Get(ctx, configurationKey(imei)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	return val, err
}

func ClearCacheConfiguration(imei string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Del(ctx, configurationKey(imei)).Err()
}
