package redis

import (
	"context"
	"neomatica/neosync/infra/constants"
	"time"
)

func deviceKey(imei string) string {
	return cacheKey("upd:imei:" + imei)
}

func SetUpdDevice(imei string) error {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.HSet(ctx, key, "update", constants.UPDATE_COMMAND).Err(); err != nil {
		return err
	}

	return Client.Expire(ctx, key, constants.Redis_DeviceUpdTTL).Err()
}

func GetUpdDevice(imei string) (string, error) {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := Client.HGet(ctx, key, constants.UPDATE_COMMAND).Result()
	if err != nil {
		return "", err
	}

	if err := Client.HDel(ctx, key, constants.UPDATE_COMMAND).Err(); err != nil {
		return "", err
	}

	return val, nil
}

func UpdDeviceExists(imei string) (bool, error) {
	key := deviceKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.HExists(ctx, key, constants.UPDATE_COMMAND).Result()
}
