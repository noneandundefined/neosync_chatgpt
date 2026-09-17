package redis

import (
	"context"
	"neomatica/neosync-tcp/infra/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

func deviceKey(imei string) string {
	return cacheKey("upd:imei:" + imei)
}

func SetUpdDevice(imei string) error {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.HSet(ctx, key, "update", constants.UPDATE_COMMAND).Err(); err != nil {
		return err
	}

	return client.Expire(ctx, key, 8*time.Minute).Err()
}

func GetUpdDevice(imei string) (*string, error) {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := client.HGet(ctx, key, constants.UPDATE_COMMAND).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := client.HDel(ctx, key, constants.UPDATE_COMMAND).Err(); err != nil {
		return nil, err
	}

	return &val, nil
}

func UpdDeviceExists(imei string) (bool, error) {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.HExists(ctx, key, constants.UPDATE_COMMAND).Result()
}

func DeleteUpdDevice(imei string) error {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}
