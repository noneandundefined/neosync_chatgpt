package redis

import (
	"context"
	"fmt"
	"neomatica/neosync/infra/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

func deviceAuthKey(userUuid string, deviceId uint64) string {
	return fmt.Sprintf("device_auth:%s:%d", userUuid, deviceId)
}

func SetDeviceAuth(userUuid string, deviceId uint64, password string) error {
	key := deviceAuthKey(userUuid, deviceId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Set(ctx, key, password, constants.Redis_DeviceAuthTTL).Err()
}

func GetDeviceAuth(userUuid string, deviceId uint64) (string, error) {
	key := deviceAuthKey(userUuid, deviceId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func DeleteDeviceAuth(userUuid string, deviceId uint64) error {
	key := deviceAuthKey(userUuid, deviceId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Del(ctx, key).Err()
}
