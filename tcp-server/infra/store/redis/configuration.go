package redis

import (
	"context"
	"time"
)

func configurationErrKey(imei string) string {
	return cacheKey("cfg-error:imei:" + imei)
}

func configurationKey(imei string) string {
	return cacheKey("cfg:imei:" + imei)
}

func SetRequireConfiguration(imei string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Set(ctx, configurationErrKey(imei), "true", 0).Err()
}

func ClearRequireConfiguration(imei string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Del(ctx, configurationErrKey(imei)).Err()
}

func ClearCacheConfiguration(imei string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Del(ctx, configurationKey(imei)).Err()
}
